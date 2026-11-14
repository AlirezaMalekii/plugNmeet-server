// speaking-observer records server-observed speech transitions independently of browser telemetry.
package main

import (
 "context"
 "encoding/json"
 "fmt"
 "log"
 "os"
 "os/signal"
 "path/filepath"
 "sync"
 "syscall"
 "time"

 "github.com/livekit/protocol/auth"
 "github.com/livekit/protocol/livekit"
 lksdk "github.com/livekit/server-sdk-go/v2"
)

type journal struct {
 mu sync.Mutex
 file *os.File
 room string
 session string
}
func (j *journal) write(kind, identity, name string) {
 j.mu.Lock()
 defer j.mu.Unlock()
 event := map[string]any{"event":kind,"room_id":j.room,"room_sid":j.session,"participant_key":identity,"name":name,"timestamp_ms":time.Now().UnixMilli()}
 if err:=json.NewEncoder(j.file).Encode(event);err!=nil {log.Printf("journal encode: %v",err);return}
 if err:=j.file.Sync();err!=nil {log.Printf("journal sync: %v",err)}
}
func observe(ctx context.Context, info *livekit.Room, output, host, key, secret string) error {
 // A SID-scoped filename separates repeated sessions using the same room name.
 f,err:=os.OpenFile(filepath.Join(output,info.Sid+".jsonl"),os.O_CREATE|os.O_APPEND|os.O_WRONLY,0600)
 if err!=nil{return err}
 defer f.Close()
 j:=journal{file:f,room:info.Name,session:info.Sid}
 done:=make(chan struct{})
 var once sync.Once
 room:=lksdk.NewRoom(&lksdk.RoomCallback{
  ParticipantCallback:lksdk.ParticipantCallback{
   OnIsSpeakingChanged:func(p lksdk.Participant){
    kind:="speech_ended"
    if p.IsSpeaking(){kind="speech_started"}
    j.write(kind,p.Identity(),p.Name())
   },
  },
  OnParticipantDisconnected:func(p *lksdk.RemoteParticipant){j.write("participant_left",p.Identity(),p.Name())},
  OnDisconnected:func(){once.Do(func(){close(done)})},
 })
 no:=false
 token,err:=auth.NewAccessToken(key,secret).SetIdentity("ineshast-observer-"+info.Sid).
  SetVideoGrant(&auth.VideoGrant{RoomJoin:true,Room:info.Name,Hidden:true,CanPublish:&no,CanSubscribe:&no}).
  SetValidFor(10*time.Minute).ToJWT()
 if err!=nil{return err}
 if err=room.JoinWithToken(host,token,lksdk.WithAutoSubscribe(false));err!=nil{return err}
 j.write("observer_connected","","")
 select {case <-ctx.Done():case <-done:}
 room.Disconnect()
 j.write("observer_disconnected","","")
 return nil
}
func main(){
 host,key,secret:=os.Getenv("LIVEKIT_URL"),os.Getenv("LIVEKIT_API_KEY"),os.Getenv("LIVEKIT_API_SECRET")
 if host==""||key==""||secret==""{log.Fatal("LIVEKIT_URL, LIVEKIT_API_KEY and LIVEKIT_API_SECRET required")}
 output:=os.Getenv("SPEAKING_OUTPUT")
 if output==""{output="/data/speaking"}
 if err:=os.MkdirAll(output,0700);err!=nil{log.Fatal(err)}
 ctx,cancel:=signal.NotifyContext(context.Background(),syscall.SIGINT,syscall.SIGTERM);defer cancel()
 client:=lksdk.NewRoomServiceClient(host,key,secret)
 var mu sync.Mutex
 active:=map[string]bool{}
 tick:=time.NewTicker(time.Second);defer tick.Stop()
 for {
  select {case <-ctx.Done():return;case <-tick.C:}
  rooms,err:=client.ListRooms(ctx,&livekit.ListRoomsRequest{})
  if err!=nil{log.Printf("list rooms: %v",err);continue}
  for _,info:=range rooms.Rooms {
   mu.Lock()
   running:=active[info.Sid]
   if !running {active[info.Sid]=true}
   mu.Unlock()
   if running {continue}
   go func(info *livekit.Room){
    defer func(){mu.Lock();delete(active,info.Sid);mu.Unlock()}()
    if err:=observe(ctx,info,output,host,key,secret);err!=nil{log.Print(fmt.Errorf("observe %s: %w",info.Sid,err))}
   }(info)
  }
 }
}
