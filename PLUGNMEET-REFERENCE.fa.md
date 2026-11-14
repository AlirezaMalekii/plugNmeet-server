# مرجع فارسی plugNmeet برای آی‌نشست

نسخهٔ مرجع در ۱۴ نوامبر ۲۰۲۶ از مخزن رسمی مستندات ثبت شده است. این فایل راهنمای فارسی، قراردادهای API و متن فنی همهٔ ۴۷ صفحهٔ نسخهٔ انگلیسی را در یک مرجع قابل جستجو گردآوری می‌کند. ۲۵ صفحه در بخش API قرار دارند. سه صفحهٔ API قدیمی analytics نیز در پوشش حفظ شده‌اند.

مبنای صفحات commit `44a474a17d965e17ebb68c43991205768ef4c27f` است. فایل‌های اصلی، دسته‌بندی‌ها و مجوز MIT در `docs/upstream/` نگهداری می‌شوند و `manifest.json` مسیر، حجم و SHA256 هر فایل را ثبت می‌کند. نسخهٔ بررسی‌شدهٔ مستندات ممکن است با نسخهٔ بعدی نرم‌افزار تفاوت داشته باشد. پیش از ارتقا، تغییر قرارداد در پروتکل و سورس همان release بررسی شود. [مخزن اصلی مستندات](https://github.com/mynaparrot/plugNmeet-docs/tree/44a474a17d965e17ebb68c43991205768ef4c27f/docs)

این سند، راهنمای پیاده‌سازی و استفاده است. وضعیت واقعی نصب، پذیرش پایلوت و انتشار در `IMPLEMENTATION-STATUS.md` ثبت می‌شود. وجود یک روش یا قابلیت در مستندات به معنای اجرا و تأیید آن روی سرور آی‌نشست نیست.

## روش خواندن و به‌روزرسانی

هر صفحه یک بخش با راهنمای فارسی، لینک متن رسمی، لینک سورس ثابت و نسخهٔ محلی دارد. پس از راهنمای فارسی، متن فنی کامل همان صفحه شامل جدول فیلدها، نمونهٔ درخواست و پاسخ، خطاهای مستند و توضیحات ویژه آمده است. این متن فنی انگلیسی برای حفظ دقیق نام و قرارداد فیلدها نگه داشته شده است. نمونه‌های رسمی دمو در این گردآوری با `EXAMPLE_API_SECRET` جایگزین شده‌اند. نسخهٔ اصلی بایگانی‌شده بدون این جایگزینی قابل مراجعه است.

برای به‌روزرسانی، ابتدا commit جدید مستندات مشخص شود، فهرست `docs/` و منوی خودکار با این ۴۷ صفحه مقایسه شود، متن صفحات و manifest تازه ساخته شوند و تغییرات فیلدها و مثال‌ها بازبینی شوند. در کنار این کار، commit سرور، کلاینت، recorder و نسخهٔ پروتکل باید ثبت شود. لینک‌هایی که به `main` در پروژه‌های وابسته اشاره دارند محتوای متغیر دارند و قبل از استفاده باید با نسخهٔ نصب‌شده تطبیق داده شوند.

## مسیر اجرایی آی‌نشست

1. Core پس از بررسی مجوز، یک جلسهٔ آنلاین را با ارائه‌دهندهٔ ذخیره‌شده ایجاد می‌کند. اتاق plugNmeet با شناسهٔ پایدار ساخته می‌شود و `enable_analytics` فعال است.
2. برای هر عضو یا مهمان، شناسهٔ پایدار و نقش تأییدشدهٔ میزبان یا شرکت‌کننده به API ورود فرستاده می‌شود. مرورگر توکن کوتاه‌عمر را دریافت می‌کند.
3. رویدادهای جلسه در گیرندهٔ امضاشدهٔ وب‌هوک دریافت و پیش از پردازش در صف پایدار ثبت می‌شوند. UUID رویداد برای جلوگیری از اثر تکراری استفاده می‌شود.
4. `session_ended` پایان اجرای جلسه را ثبت می‌کند. `recording_proceeded` فایل ضبط آماده و `analytics_proceeded` گزارش آماده را اعلام می‌کنند. ترتیب رسیدن این سه رویداد مستقل است.
5. فایل‌های ضبط و JSON گزارش با API و توکن موقت دریافت می‌شوند. هر قطعهٔ ضبط به صورت مستقل به جلسه متصل می‌شود و فایل صوتی قابل استفاده برای پردازش گفتار تولید می‌شود.
6. برای بازهٔ گفتار، وضعیت گویندگان سمت سرور ثبت و با شروع هر قطعهٔ ضبط هم‌زمان می‌شود. `talked_duration` گزارش بومی برای تطبیق مجموع‌ها استفاده می‌شود. زمان میکروفن باز با زمان صحبت کاربر تفاوت دارد.
7. کار زمان‌بندی‌شده از فهرست نشست‌ها، ضبط‌ها و artifactها خروجی جامانده را بازیابی می‌کند. خطای موقت وب‌هوک نباید باعث از دست رفتن خروجی قابل دریافت شود.

## قرارداد احراز درخواست‌های مدیریتی

پایهٔ API `https://meet.example.com/auth` است. بیشتر endpointهای فهرست‌شده با `POST` و JSON فراخوانی می‌شوند. مثال زیر فقط نمونه است و اعتبار واقعی باید از محیط امن backend خوانده شود.

```php
$body = json_encode(['room_id' => 'room01'], JSON_UNESCAPED_SLASHES);
$signature = hash_hmac('sha256', $body, $apiSecret);
// Send these exact $body bytes with API-KEY and HASH-SIGNATURE headers.
```

هدرهای لازم `API-KEY`، `HASH-SIGNATURE` و `Content-Type: application/json` هستند. امضا از بایت‌های دقیق بدنه ساخته می‌شود. برای `room/uploadWhiteboardFile` بدنه multipart است و امضا از مقدار `Room-Id` ساخته می‌شود. جدول پاسخ هر endpoint، `status` و `status_code` کاربردی را مشخص می‌کند. HTTP موفق به تنهایی موفقیت کاربردی را اثبات نمی‌کند.

## عملیات، نگهداری و پذیرش

منبع سرور در فورک `AlirezaMalekii/plugNmeet-server` نگهداری می‌شود. `origin` به فورک و `upstream` به پروژهٔ رسمی اشاره دارد. کلاینت و recorder در submoduleها با commit ثابت قرار می‌گیرند. کد اتصال Laravel در مخزن آی‌نشست است. انتشار باید SHA کد و digest تصویر را ثبت کند و فایل پیکربندی واقعی و اسرار خارج از مخزن عمومی نگهداری شوند.

دامنه‌های هدف `meet.ineshast.ir` و `turn.ineshast.ir` هستند. در معماری میزبان مشترک، ورودی HAProxy مالک پورت‌های HTTP/HTTPS است و نام TURN به backend مخصوص خود هدایت می‌شود. شبکه و volume و اعتبارنامهٔ plugNmeet مستقل از سامانه‌های دیگر تعریف می‌شوند. پورت‌های داخلی دیتابیس، Redis و NATS با نیاز ارتباطی محدود در دسترس قرار می‌گیرند. ظرفیت recorder با اندازه‌گیری بار واقعی تعیین می‌شود.

برای ساخت و انتشار، دستورهای دقیق باید از فایل‌های `ops/` همان commit اجرا شوند. خروجی build، digest تصویر، SHA هر submodule و نسخهٔ Compose ثبت شود. پیش از جابه‌جایی نسخه، پشتیبان سازگار دیتابیس، پیکربندی و فایل‌های ماندگار تهیه شود. بازگشت از digest قبلی اجرا می‌شود و بازگشت دیتابیس فقط بر اساس سازگاری migration یا پشتیبان آزموده‌شده انجام شود. نمونهٔ انتشار رسمی در بخش نصب آمده است و باید با تنظیمات این میزبان تطبیق داده شود.

تمدید گواهی باید با dry-run و بازخوانی گواهی سرویس ورودی آزموده شود. پشتیبان باید شامل دیتابیس، تنظیمات امن و فایل‌های ضبط و artifact باشد و رمزگذاری و یک بازیابی مستقل داشته باشد. وضعیت backup و فضای دیسک و مدت انتظار صف ضبط پایش شود. سیاست حذف محلی باید پس از دریافت و تأیید خروجی با سیاست نگهداری آی‌نشست هماهنگ باشد.

اگر کاربر وارد می‌شود و صدا ندارد، مسیر UDP، IP عمومی اعلام‌شدهٔ LiveKit، پورت fallback و TURN بررسی شوند. اگر رابط جلسه همگام نمی‌شود، اتصال WebSocket به NATS و پیکربندی دامنه بررسی شود. اگر فایل آماده نمی‌شود، وضعیت recorder، فضای دیسک، صف و hookهای پردازش بررسی شوند. اگر فایل آماده است و در آی‌نشست دیده نمی‌شود، امضای وب‌هوک، نگاشت room/sid، صف Core و reconciliation بررسی شوند.

پذیرش پایلوت شامل ۲ اتاق، مجموعاً ۱۰ شرکت‌کننده و ۲ ضبط هم‌زمان است. از دو شبکهٔ بیرونی، صدا و تصویر و اشتراک صفحه و اتصال اجباری TURN آزموده می‌شوند. گفتار معلوم با سکوت، هم‌پوشانی، قطع اتصال و چند قطعهٔ ضبط باید هویت صحیح و مرزهای زمانی با خطای حداکثر ۲ ثانیه بدهد. رویداد تکراری و دیرهنگام، قطع موقت Core، بازیابی فایل از API، پخش ضبط قدیمی BBB و بازگشت انتشار نیز جزو شواهد پذیرش هستند. پیش‌فرض جلسه‌های جدید پس از قبولی این آزمون‌ها تغییر می‌کند.

## نمایهٔ همهٔ صفحات

| صفحه | موضوع فارسی | بخش |
| --- | --- | --- |
| [api/analytics/delete](#page-01) | حذف analytics قدیمی | API |
| [api/analytics/download](#page-02) | دانلود analytics قدیمی | API |
| [api/analytics/fetch](#page-03) | فهرست analytics قدیمی | API |
| [api/artifact/artifact-info](#page-04) | جزئیات artifact | API |
| [api/artifact/delete](#page-05) | حذف artifact | API |
| [api/artifact/download](#page-06) | توکن دانلود artifact | API |
| [api/artifact/fetch](#page-07) | فهرست artifactهای جلسه | API |
| [api/get-client-files](#page-08) | دریافت فایل‌های رابط جلسه | API |
| [api/intro](#page-09) | اصول احراز هویت API | API |
| [api/recording/delete](#page-10) | حذف ضبط | API |
| [api/recording/download](#page-11) | توکن دانلود ضبط | API |
| [api/recording/fetch](#page-12) | فهرست فایل‌های ضبط | API |
| [api/recording/merge-recordings](#page-13) | ادغام ضبط‌های چندبخشی | API |
| [api/recording/recording-info](#page-14) | جزئیات یک فایل ضبط | API |
| [api/recording/update-metadata](#page-15) | به‌روزرسانی metadata ضبط | API |
| [api/room/broadcast-to-room](#page-16) | ارسال پیام و اعلان به اتاق | API |
| [api/room/create-poll](#page-17) | ساخت نظرسنجی و آزمون | API |
| [api/room/create](#page-18) | ساخت اتاق و تنظیم قابلیت‌ها | API |
| [api/room/end](#page-19) | پایان اتاق | API |
| [api/room/fetch-past](#page-20) | فهرست نشست‌های گذشته | API |
| [api/room/is-active](#page-21) | بررسی فعال‌بودن اتاق | API |
| [api/room/join](#page-22) | توکن ورود و هویت کاربر | API |
| [api/room/room-info](#page-23) | جزئیات اتاق فعال | API |
| [api/room/rooms-info](#page-24) | فهرست اتاق‌های فعال | API |
| [api/room/upload-whiteboard-file](#page-25) | ارسال فایل به وایت‌برد زنده | API |
| [benchmark](#page-26) | ظرفیت و روش آزمون | راهنما |
| [developer-guide/design-customisation](#page-27) | شخصی‌سازی ظاهر | راهنما |
| [developer-guide/intro](#page-28) | مسیر توسعهٔ اتصال اختصاصی | راهنما |
| [developer-guide/mobile-app-integration](#page-29) | اتصال اپلیکیشن موبایل و حالت Hybrid | راهنما |
| [developer-guide/scalable-setup](#page-30) | مقیاس‌پذیری و جداسازی ضبط‌کننده | راهنما |
| [developer-guide/setup-development](#page-31) | محیط توسعهٔ محلی | راهنما |
| [firewall](#page-32) | فایروال و پورت‌های رسانه | راهنما |
| [installation](#page-33) | نصب و نگهداری پایه | راهنما |
| [intro](#page-34) | معرفی و معماری ابزار | راهنما |
| [others/analytics](#page-35) | گزارش حضور، مشارکت و زمان گفتار | راهنما |
| [others/hooks](#page-36) | Hookهای اسکریپت و ذخیره‌سازی | راهنما |
| [others/webhooks](#page-37) | وب‌هوک‌ها و چرخهٔ جلسه | راهنما |
| [security-overview](#page-38) | امنیت، حریم خصوصی و چرخهٔ داده | راهنما |
| [tutorials/migration-from-bbb](#page-39) | سازگاری و مهاجرت از BBB | راهنما |
| [tutorials/quick_php](#page-40) | نمونهٔ PHP برای ایجاد و ورود | راهنما |
| [user-guide/attendee](#page-41) | راهنمای شرکت‌کننده | راهنما |
| [user-guide/joomla-integration](#page-42) | اتصال Joomla | راهنما |
| [user-guide/lti](#page-43) | اتصال آموزشی با LTI | راهنما |
| [user-guide/moderator](#page-44) | راهنمای مدیر جلسه | راهنما |
| [user-guide/moodle-integration](#page-45) | اتصال Moodle | راهنما |
| [user-guide/overview](#page-46) | مسیر راهنمای کاربران | راهنما |
| [user-guide/wordPress-integration](#page-47) | اتصال WordPress | راهنما |

## فهرست endpointهای API

مسیرهای جدول به `/auth` افزوده می‌شوند. دانلود نهایی با توکن مسیر GET جداگانه دارد که در صفحهٔ مربوط آمده است.

| صفحه | مسیر دقیق | وضعیت |
| --- | --- | --- |
| حذف analytics قدیمی | `/analytics/delete` | منسوخ، استفاده از artifact برای اتصال جدید |
| دانلود analytics قدیمی | `/analytics/getDownloadToken` | منسوخ، استفاده از artifact برای اتصال جدید |
| فهرست analytics قدیمی | `/analytics/fetch` | منسوخ، استفاده از artifact برای اتصال جدید |
| جزئیات artifact | `/artifact/info` | مستند در نسخهٔ مرجع |
| حذف artifact | `/artifact/delete` | مستند در نسخهٔ مرجع |
| توکن دانلود artifact | `/artifact/getDownloadToken` | مستند در نسخهٔ مرجع |
| فهرست artifactهای جلسه | `/artifact/fetch` | مستند در نسخهٔ مرجع |
| دریافت فایل‌های رابط جلسه | `/getClientFiles` | مستند در نسخهٔ مرجع |
| حذف ضبط | `/recording/delete` | مستند در نسخهٔ مرجع |
| توکن دانلود ضبط | `/recording/getDownloadToken` | مستند در نسخهٔ مرجع |
| فهرست فایل‌های ضبط | `/recording/fetch` | مستند در نسخهٔ مرجع |
| ادغام ضبط‌های چندبخشی | `/recording/mergeRecordings` | مستند در نسخهٔ مرجع |
| جزئیات یک فایل ضبط | `/recording/info` | مستند در نسخهٔ مرجع |
| به‌روزرسانی metadata ضبط | `/recording/updateMetadata` | مستند در نسخهٔ مرجع |
| ارسال پیام و اعلان به اتاق | `/room/broadcastToRoom` | مستند در نسخهٔ مرجع |
| ساخت نظرسنجی و آزمون | `/room/createPoll` | مستند در نسخهٔ مرجع |
| ساخت اتاق و تنظیم قابلیت‌ها | `/room/create` | مستند در نسخهٔ مرجع |
| پایان اتاق | `/room/endRoom` | مستند در نسخهٔ مرجع |
| فهرست نشست‌های گذشته | `/room/fetchPastRooms` | مستند در نسخهٔ مرجع |
| بررسی فعال‌بودن اتاق | `/room/isRoomActive` | مستند در نسخهٔ مرجع |
| توکن ورود و هویت کاربر | `/room/getJoinToken` | مستند در نسخهٔ مرجع |
| جزئیات اتاق فعال | `/room/getActiveRoomInfo` | مستند در نسخهٔ مرجع |
| فهرست اتاق‌های فعال | `/room/getActiveRoomsInfo` | مستند در نسخهٔ مرجع |
| ارسال فایل به وایت‌برد زنده | `/room/uploadWhiteboardFile` | مستند در نسخهٔ مرجع |

<a id="page-01"></a>

## حذف analytics قدیمی

حذف از خانوادهٔ analytics منسوخ در این بخش مستند شده است. در اتصال جدید چرخهٔ نگهداری artifact و endpoint حذف آن مبنا قرار گیرد. برای سازگاری، نام فیلد و پاسخ قدیمی در پیوست حفظ شده‌اند.

این بخش صفحهٔ `api/analytics/delete` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/analytics/delete) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/analytics/delete.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/analytics/delete.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Delete analytics

**Deprecated**
This endpoint is deprecated and will be removed in a future version. Please use the new, more flexible **[`/artifact/delete`](https://www.plugnmeet.org/docs/api/artifact/delete)** endpoint instead.

You can delete an analytics artifact by using its `artifact_id` from the `/artifact/fetch` response.


Endpoint: `/analytics/delete`

| Field     | Type   | Position | Required | Description                          |
| --------- | ------ | -------- | :------- | ------------------------------------ |
| file_id | string | root     | Yes      | File Id that you want to delete |

**Example**:

```json
{
  "file_id": "RM_SqZLoUieWYu8-209-1693729854"
}
```

##### Response

| Field  | Type    | Position | Description               |
| :----- | ------- | -------- | :------------------------ |
| status | boolean | root     | The status of the request |
| msg    | string  | root     | Response message          |

<a id="page-02"></a>

## دانلود analytics قدیمی

این endpoint منسوخ توکن دانلود analytics را صادر می‌کند. اتصال جدید باید artifact/getDownloadToken و download/artifact را استفاده کند. توکن قدیمی نیز لینک دائمی فایل نیست.

این بخش صفحهٔ `api/analytics/download` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/analytics/download) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/analytics/download.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/analytics/download.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Download Analytics (Token)

**Deprecated**
This endpoint is deprecated and will be removed in a future version. Please use the new, more flexible **[`/artifact/getDownloadToken`](https://www.plugnmeet.org/docs/api/artifact/download)** endpoint instead.

You can get a download token for an analytics artifact by using its `artifact_id` from the `/artifact/fetch` response.


Endpoint: `/analytics/getDownloadToken`

| Field   | Type   | Position | Required | Description                       |
| ------- | ------ | -------- | :------- | --------------------------------- |
| file_id | string | root     | Yes      | The ID of the file you want to download. |

**Example Request:**

```json
{
  "file_id": "RM_SqZLoUieWYu8-209-1693729854"
}
```

##### Response

| Field  | Type    | Position | Description                |
| :----- | ------- | -------- | :------------------------- |
| status | boolean | root     | Indicates if the request was successful. |
| msg    | string  | root     | Response message.          |
| token  | string  | root     | Download token.            |

Once you receive the token, you can construct the download URL using the following format:
```
https://Your-Plug-N-Meet-Server.com/download/analytics/<TOKEN HERE>
```

<a id="page-03"></a>

## فهرست analytics قدیمی

این خانوادهٔ endpoint منسوخ است. برای اتصال جدید، artifact/fetch با type برابر MEETING_ANALYTICS استفاده شود. قرارداد قدیمی برای خواندن داده‌ها و پشتیبانی اتصال‌های پیشین در پیوست حفظ شده است.

این بخش صفحهٔ `api/analytics/fetch` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/analytics/fetch) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/analytics/fetch.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/analytics/fetch.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Fetch Analytics

**Deprecated**
This endpoint is deprecated and will be removed in a future version. Please use the new, more flexible **[`/artifact/fetch`](https://www.plugnmeet.org/docs/api/artifact/fetch)** endpoint instead.

To fetch analytics data, you can make the following request:
`POST /auth/artifact/fetch`
With the body: `{"type": "MEETING_ANALYTICS"}`


Endpoint: `/analytics/fetch`

##### Request Parameters

| Field    | Type   | Required | Description                                   |
| -------- | ------ | -------- | --------------------------------------------- |
| room_ids | array  | Yes      | Array of room IDs to fetch analytics for.     |
| from     | number | No       | Starting index for records. Default is 0.     |
| limit    | number | No       | Maximum number of records to return. Default is 20. |
| order_by | string | No       | Sort order: `DESC` or `ASC`. Default is `DESC`. |

**Example Request:**

```json
{
  "room_ids": ["room01"],
  "from": 0,
  "limit": 20,
  "order_by": "DESC"
}
```

##### Response

| Field             | Type                        | Position | Description                       |
| ----------------- | -------------------------- | -------- | --------------------------------- |
| status            | boolean                     | root     | Indicates if the request was successful. |
| msg               | string                      | root     | Response message.                 |
| [result](https://www.plugnmeet.org/docs/api/analytics/fetch#result) | object                      | root     | Contains the analytics data.      |

###### Result

| Field            | Type                                      | Description                                 |
| ---------------- | ----------------------------------------- | ------------------------------------------- |
| total_analytics  | number                                    | Total number of analytics records found.    |
| from             | number                                    | Starting index for the returned records.    |
| limit            | number                                    | Number of records returned.                 |
| order_by         | string                                    | Sort order used for the records.            |
| analytics_list   | Array\<[analytics-info](https://www.plugnmeet.org/docs/api/analytics/fetch#analytics-info)> | List of analytics records.                  |

###### Analytics Info

| Field              | Type   | Description                                 |
| ------------------ | ------ | ------------------------------------------- |
| room_id            | string | The ID of the room.                         |
| file_id            | string | The ID of the analytics file.               |
| file_name          | string | The name of the analytics file.             |
| file_size          | number | The size of the file in bytes.              |
| creation_time      | number | File creation time (Unix timestamp).        |
| room_creation_time | number | Room creation time (Unix timestamp).        |

<a id="page-04"></a>

## جزئیات artifact

اطلاعات یک artifact و metadata آن با شناسهٔ مربوط دریافت می‌شود. نوع خروجی باید پیش از تفسیر محتوا بررسی شود. شکل ArtifactInfo و RoomArtifactMetadata در جدول‌های پیوست آمده است.

این بخش صفحهٔ `api/artifact/artifact-info` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/artifact/artifact-info) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/artifact/artifact-info.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/artifact/artifact-info.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Get Artifact Info

Endpoint: `/artifact/info`

This endpoint allows you to retrieve information about a single artifact using its unique `artifact_id`.

##### Request Parameters

| Field       | Type   | Required | Description                            |
| ----------- | ------ | -------- | -------------------------------------- |
| artifact_id | string | Yes      | The unique identifier of the artifact. |

###### Example Request

```json
{
  "artifact_id": "ARTIFACT_ID"
}
```

##### Response

| Field         | Type         | Description                                                                                                                    |
| :------------ |:-------------|--------------------------------------------------------------------------------------------------------------------------------|
| status        | boolean      | Indicates if the request was successful.                                                                                       |
| msg           | string       | Response message.                                                                                                              |
| status_code | string       | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |
| artifact_info | ArtifactInfo | Details about the requested artifact.                                                                                          |
| room_info     | PastRoomInfo | Details about the room session this artifact belongs to.                                                                       |

###### ArtifactInfo

| Field       | Type               | Description                                                     |
| ----------- | ------------------ | --------------------------------------------------------------- |
| artifact_id | string             | The unique identifier for this artifact.                        |
| room_id     | string             | The ID of the room this artifact belongs to.                    |
| type        | string             | The type of the artifact. See [Artifact Types](https://www.plugnmeet.org/docs/api/artifact/fetch#artifact-types). |
| created     | string             | The creation timestamp of the artifact.                         |
| metadata    | RoomArtifactMetadata | Additional metadata about the artifact, which varies by type.   |

###### RoomArtifactMetadata

This object contains metadata specific to the artifact's type. For example, for a `MEETING_SUMMARY_USAGE` artifact, it might contain `token_usage` and `cost` information. The exact structure depends on the artifact type.

<a id="page-05"></a>

## حذف artifact

حذف artifact یک عملیات مدیریتی روی خروجی ذخیره‌شده است. پیش از حذف، ارتباط آن با تحلیل، گزارش و نگهداری جلسه بررسی شود و نتیجهٔ API در وضعیت محلی منعکس شود.

این بخش صفحهٔ `api/artifact/delete` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/artifact/delete) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/artifact/delete.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/artifact/delete.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Delete Artifact

Endpoint: `/artifact/delete`

This endpoint allows you to delete an artifact file from the server's filesystem.

**Which Artifacts Are Deletable?**
This endpoint is only for deleting the physical files associated with certain artifact types (e.g., `MEETING_SUMMARY`, `SPEECH_TRANSCRIPTION`). It does not apply to metadata-only artifacts.


**Data Retention**
This action deletes the artifact **file** (e.g., the VTT or summary text file) to save storage space. The associated **metadata** about the artifact (e.g., its ID, type, and usage stats) is retained in the database for historical and auditing purposes.


##### Request Parameters

| Field       | Type   | Required | Description                            |
| ----------- | ------ | -------- | -------------------------------------- |
| artifact_id | string | Yes      | The unique identifier of the artifact. |

###### Example Request

```json
{
  "artifact_id": "ARTIFACT_ID"
}
```

##### Response

| Field  | Type    | Description                              |
| :----- | :------ | ---------------------------------------- |
| status | boolean | Indicates if the request was successful. |
| msg    | string  | Response message.                        |
| status_code | string  | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |

<a id="page-06"></a>

## توکن دانلود artifact

برای فایل artifact توکن موقت گرفته می‌شود و مسیر GET /download/artifact/<TOKEN> استفاده می‌شود. دریافت MEETING_ANALYTICS و رونویسی باید از این مسیر نسخهٔ جدید انجام شود. اصل JSON دریافت‌شده برای امکان پردازش مجدد حفظ شود.

این بخش صفحهٔ `api/artifact/download` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/artifact/download) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/artifact/download.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/artifact/download.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Get Download Token

Endpoint: `/artifact/getDownloadToken`

This endpoint generates a secure, one-time-use token that can be used to download an artifact file.

**Which Artifacts Are Downloadable?**
This endpoint is only for artifact types that represent a physical file, such as `MEETING_SUMMARY` or `SPEECH_TRANSCRIPTION`. It will not work for metadata-only artifacts like `MEETING_SUMMARY_USAGE`.


##### Request Parameters

| Field       | Type   | Required | Description                            |
| ----------- | ------ | -------- | -------------------------------------- |
| artifact_id | string | Yes      | The unique identifier of the artifact. |

###### Example Request

```json
{
  "artifact_id": "ARTIFACT_ID"
}
```

##### Response

| Field  | Type   | Description                      |
| :----- | :----- | -------------------------------- |
| status | boolean| Indicates if the request was successful. |
| msg    | string | Response message.                |
| status_code | string  | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |
| token  | string | The one-time-use download token.   |

After receiving the token, you can provide the following download URL to your user:
`https://Your-Plug-N-Meet-Server.com/download/artifact/THE_TOKEN_HERE`

**Downloading meeting analytics?**
For the full workflow — list `MEETING_ANALYTICS` artifacts, download the JSON, and format it into attendance and engagement reports with the PHP SDK — see the [Meeting Analytics & Reports guide](https://www.plugnmeet.org/docs/others/analytics).


<a id="page-07"></a>

## فهرست artifactهای جلسه

artifact خروجی‌هایی مانند MEETING_ANALYTICS، MEETING_SUMMARY و SPEECH_TRANSCRIPTION را در بر می‌گیرد. room_ids، room_sid و type فیلترهای اصلی‌اند. صفحه‌بندی from و limit دارد و سقف limit این نسخه ۱۰۰ است. برای گزارش جلسهٔ مشخص، room_sid را همراه room_id استفاده کنید.

این بخش صفحهٔ `api/artifact/fetch` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/artifact/fetch) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/artifact/fetch.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/artifact/fetch.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Fetch Artifacts

Endpoint: `/artifact/fetch`

This endpoint allows you to retrieve a paginated list of all generated meeting artifacts. You can filter the results by room, session, or artifact type.

##### Request Parameters

| Field      | Type             | Required | Description                                                                                                                                                             |
| ---------- | ---------------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| room_ids   | array of strings | No       | An array of one or more `room_id`s to filter the results.                                                                                                               |
| room_sid   | string           | No       | Filter artifacts from a specific room session by providing its `sid`.                                                                                                   |
| type       | string or number | No       | Filter by a specific artifact type. You can use either the string name (e.g., `MEETING_ANALYTICS`) or its corresponding integer value. See Artifact Types below. |
| from       | number           | No       | The starting offset for pagination. Default: `0`.                                                                                                                       |
| limit      | number           | No       | The maximum number of artifacts to return. Default: `20`, Maximum: `100`.                                                                                               |
| order_by   | string           | No       | The sort order for the results. Can be `ASC` or `DESC`. Default: `DESC`.                                                                                                |

###### Example Request

```json
{
  "room_ids": ["room01", "room02"],
  "type": "MEETING_ANALYTICS",
  "limit": 10
}
```

##### Response

| Field  | Type              | Description                                |
| :----- |-------------------| ------------------------------------------ |
| status | boolean           | Indicates if the request was successful.   |
| msg    | string            | Response message.                          |
| status_code | string  | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |
| result | [Result](https://www.plugnmeet.org/docs/api/artifact/fetch#result) | The object containing the list of artifacts. |

###### Result

| Field           | Type                  | Description                                                                                             |
| --------------- | --------------------- | ------------------------------------------------------------------------------------------------------- |
| total_artifacts | number                | The total number of artifacts matching the query.                                                       |
| from            | number                | The pagination offset of the current result set.                                                        |
| limit           | number                | The pagination limit of the current result set.                                                         |
| order_by        | string                | The sorting order of the result set.                                                                    |
| type            | string                | The artifact type that was filtered for.                                                                |
| artifacts_list  | array of ArtifactInfo | The list of artifact objects. See the **[Get Artifact Info](https://www.plugnmeet.org/docs/api/artifact/artifact-info)** page for the detailed structure of the `ArtifactInfo` object. |

##### Artifact Types

The `type` field identifies the kind of artifact. You can use either the string name (e.g., `MEETING_ANALYTICS`) or its corresponding integer value when filtering.

A few common examples include:
*   `MEETING_SUMMARY`
*   `SPEECH_TRANSCRIPTION`
*   `MEETING_ANALYTICS`
*   `MEETING_SUMMARY_USAGE`

For a complete and up-to-date list of all available artifact types, please refer to the official [**`plugnmeet_room_artifacts.proto`**](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_room_artifacts.proto) file.

**Working with meeting analytics?**
To list and download `MEETING_ANALYTICS` files and turn the raw JSON into attendance, engagement, and classroom performance reports, see the [Meeting Analytics & Reports guide](https://www.plugnmeet.org/docs/others/analytics).


<a id="page-08"></a>

## دریافت فایل‌های رابط جلسه

با بدنهٔ خالی فهرست css_files و js_files گرفته می‌شود. نام فایل‌ها با مسیر assets و زیرپوشهٔ css یا js کامل می‌شود و div با id برابر plugNmeet-app محل اتصال React است. window.plugNmeetConfig و متغیرهای لازم پیش از scriptها تنظیم می‌شوند. نام آرایه‌های توضیحی متن و نام فیلدهای جدول با هم مقایسه شوند و پاسخ واقعی نسخهٔ استقرار مبنا باشد.

این بخش صفحهٔ `api/get-client-files` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/get-client-files) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/get-client-files.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/get-client-files.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Get Client Files

Endpoint: `/getClientFiles`

The Plug-N-Meet client is a powerful, standalone React application. While you can embed it using an iframe, we **strongly recommend** a more seamless approach using this API endpoint.

The `getClientFiles` API provides you with a list of all the necessary CSS and JavaScript files to render the Plug-N-Meet interface directly within your own web page. This method offers deeper integration, better performance, and avoids the complexities of iframes, allowing you to create a truly native-feeling video conferencing experience for your users.

It also gives you complete branding freedom. Since this method isn't tied to a specific domain, you can host the client on any SSL-enabled domain or subdomain, ensuring a seamless brand experience for your users.

###### How to Use

1.  Call the `/getClientFiles` endpoint from your server-side application.
2.  The API will respond with two arrays: `css` and `js`, containing the URLs to the required asset files.
3.  Construct the full URL for each filename. The API returns plain filenames (e.g., `app.js`). You must prepend your Plug-N-Meet server's asset path and the correct subdirectory (`/js/` for JavaScript files, `/css/` for CSS files). For example: `https://plugnmeet.example.com/assets/js/app.js`.
4.  In the HTML page where you want to display the conference, dynamically generate `<link>` and `<script>` tags using the full URLs you just created.
5.  Ensure you have a `div` with the id `plugNmeet-app` in your `<body>`, as this is where the React application will mount itself.

This process ensures that your application always loads the correct and most up-to-date client assets.

For implementation examples, refer to the [conference.php](https://github.com/mynaparrot/plugNmeet-sdk-php/blob/main/examples/conference.php) file or our plugins for [Moodle](https://github.com/mynaparrot/moodle-mod_plugnmeet/blob/main/conference.php), [WordPress](https://github.com/mynaparrot/plugNmeet-WordPress/blob/main/plugnmeet/public/partials/plugnmeet-public-display-client.php), and [Joomla](https://github.com/mynaparrot/plugNmeet-Joomla/blob/main/components/com_plugnmeet/tmpl/room/conference.php).


**Example Request:**

```json
{}
```

##### Response

| Field       | Type    | Position | Description                                                                                                                    |
|:------------|:--------| :------- |:-------------------------------------------------------------------------------------------------------------------------------|
| status      | boolean | root     | Indicates whether the request was successful                                                                                   |
| msg         | string  | root     | Response message                                                                                                               |
| status_code | string  | root     | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |
| css_files         | array   | root     | An array of required CSS filenames.                                                                                            |
| js_files          | array   | root     | An array of required JavaScript filenames.                                                                                     |

###### Example HTML Structure

Here is a simplified example using PHP to illustrate the logic. This approach can be easily adapted to any server-side language (like Node.js, Python, or Ruby) to dynamically generate the required HTML.

```php
<?php
// 1. Define your server URL.
$plugnmeet_server_url = 'https://plugnmeet.example.com';

// 2. Call the /getClientFiles API and get the response.
// This is a sample response for demonstration.
$api_response = [
    'status' => true,
    'css_files' => ['styles.8f34.css', 'vendor.9c3a.css'],
    'js_files' => ['runtime.3e4a.js', 'vendor.b12c.js', 'app.5d6e.js', 'main-module.a4f1.js'],
];

// 3. Construct the base path for your assets.
$assets_path = $plugnmeet_server_url . '/assets';

?>
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8"/>
    <meta http-equiv="X-UA-Compatible" content="IE=edge"/>
    <meta name="viewport" content="width=device-width,initial-scale=1"/>
    <title>plugNmeet</title>

    <!-- 4. Dynamically generate <link> tags for CSS files. -->
    <?php foreach ($api_response['css_files'] as $css_file): ?>
        <link href="<?php echo $assets_path . '/css/' . $css_file ?>" rel="stylesheet" />
    <?php endforeach; ?>

    <!--
    5. Inject required window variables before loading scripts.
       These are essential for the client to initialize correctly.
    -->
    <script type="text/javascript">
        window.plugNmeetConfig = {
            // Required: The URL of your plugNmeet server.
            serverUrl: "<?php echo $plugnmeet_server_url ?>",

            // Required: The public path to the assets directory.
            staticAssetsPath: "<?php echo $assets_path ?>",

            // Optional: Add any other custom configurations.
            // See: https://github.com/mynaparrot/plugNmeet-client/blob/main/src/assets/config_sample.js
            enableSimulcast: true,
            enableDynacast: true
        };
    </script>
</head>
<body>
    <!-- This is the root element for the Plug-N-Meet client -->
    <div id="plugNmeet-app"></div>

    <!-- 6. Dynamically generate <script> tags for JS files. -->
    <?php foreach ($api_response['js_files'] as $js_file): ?>
        <?php if (str_starts_with($js_file, 'main-module.')): ?>
            <script src="<?php echo $assets_path . '/js/' . $js_file ?>" type="module"></script>
        <?php else: ?>
            <script src="<?php echo $assets_path . '/js/' . $js_file ?>" defer="defer"></script>
        <?php endif; ?>
    <?php endforeach; ?>
</body>
</html>
```

<a id="page-09"></a>

## اصول احراز هویت API

درخواست‌های مدیریتی JSON با POST به https://meet.example.com/auth و مسیر endpoint ارسال می‌شوند. API-KEY، HASH-SIGNATURE و Content-Type لازم‌اند. امضا، HMAC-SHA256 از بایت‌های دقیق بدنه با secret و خروجی hex کوچک است. رشته‌ای که امضا می‌شود باید همان رشتهٔ ارسال‌شده باشد. آپلود multipart وایت‌برد قرارداد ویژهٔ خود را دارد. وضعیت کاربردی پاسخ در status و status_code بررسی شود و secret در کلاینت قرار نگیرد.

این بخش صفحهٔ `api/intro` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/intro) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/intro.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/intro.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Introduction

Plug-N-Meet provides a secure API for interacting with your server. All authenticated API requests share a common base path.

##### Live Demo

Explore Plug-N-Meet's features and test its API with our live demo.

**Explore the features:** [https://demo.plugnmeet.com/landing.html](https://demo.plugnmeet.com/landing.html)

**Demo API Credentials for Developers:**
```
plugNmeet server URL: https://demo.plugnmeet.com
plugNmeet API KEY: plugnmeet
plugNmeet API SECRET: EXAMPLE_API_SECRET
```

###### Endpoint Structure

To use the API, you will send a `POST` request to an endpoint that starts with your server's URL followed by the `/auth` path. The specific API method path is then appended to it.

**Base authenticated API path:**
```
https://plugnmeet.example.com/auth
```

For example, if you are calling the `/room/getActiveRoomInfo` method, the full endpoint URL would be:
```
https://plugnmeet.example.com/auth/room/getActiveRoomInfo
```

All API calls must be `POST` requests with a JSON body.

###### Authentication Headers

All requests to these endpoints must include the following headers for authentication:

- **`API-KEY`**: Your Plug-N-Meet Server API Key
- **`HASH-SIGNATURE`**: An HMAC SHA256 signature (in hex format) of the raw JSON request body, generated using your API Secret as the key.
- **`Content-Type`**: `application/json`

---

##### Generating the HASH-SIGNATURE

The `HASH-SIGNATURE` is the core of the authentication process. It proves that the request is coming from a trusted source and that the message body has not been tampered with. It is a standard HMAC (Hash-based Message Authentication Code) signature.

Here is the language-agnostic process to generate it:

1.  **Get the Raw Request Body**: Start with the complete JSON request body as a single, raw string. It is crucial that this string is exactly what will be sent in the `POST` request.
2.  **Select the Algorithm**: Use the **HMAC-SHA256** algorithm. This is a standard function available in the cryptography or hashing library of almost any programming language.
3.  **Use Your Secret Key**: The "key" for the HMAC function is your **API Secret**.
4.  **Compute the Hash**: Pass the raw request body string and your API Secret to the HMAC-SHA256 function.
5.  **Encode the Result**: The output of the HMAC function is a binary hash. You must convert this binary value into its **lowercase hexadecimal representation**. This final hex string is the value for your `HASH-SIGNATURE` header.

By following these steps, you can generate a valid signature in any language, ensuring secure communication with the API.

##### API Categories

The Plug-N-Meet API is organized into several categories based on functionality.

*   **Room Management:** Endpoints for creating, managing, and querying the state of live meeting rooms.
*   **Authentication:** Endpoints for generating access tokens for participants.
*   **Recording Management:** Endpoints for fetching, downloading, and deleting meeting recordings.
*   **Artifact Management:** Endpoints for fetching, downloading, and deleting AI-generated meeting artifacts like summaries and transcriptions.
*   **Webhooks & Scripting:** Automate workflows and integrate external services using real-time [webhook events](https://www.plugnmeet.org/docs/others/webhooks) or by implementing custom [scripting and storage hooks](https://www.plugnmeet.org/docs/others/hooks).

---

##### Implementation Examples

The following examples demonstrate how to implement the signature generation process in various popular languages.

###### Shell with cURL

This example demonstrates the entire process using command-line tools. It shows how to construct the headers and body, generate the signature, and send the request.

```bash
# 1. The raw JSON body as a string
BODY='{"room_id":"room01"}'

# Your API credentials
API_KEY="plugnmeet"
SECRET="EXAMPLE_API_SECRET"

# 2. Generate the HMAC-SHA256 signature and encode it in hex
#    (This command performs steps 2-5 from the process above)
SIGNATURE=$(echo -n "$BODY" | openssl dgst -sha256 -mac HMAC -macopt key:"$SECRET" | awk '{print $2}')

# 3. Make the POST request with the correct headers and body
curl -X POST https://demo.plugnmeet.com/auth/room/getActiveRoomInfo \
  -H "Content-Type: application/json" \
  -H "API-KEY: $API_KEY" \
  -H "HASH-SIGNATURE: $SIGNATURE" \
  -d "$BODY"
```

###### PHP

```php
$body = json_encode([
    "room_id" => "room01",
]);
$signature = hash_hmac("sha256", $body, "SECRET");
```

###### NodeJS

```js
// using CryptoJS library
const body = JSON.stringify({ room_id: "room01" });
const hash = CryptoJS.HmacSHA256(body, "SECRET");
const signature = CryptoJS.enc.Hex.stringify(hash);
```

###### TypeScript (Node.js built-in crypto)

```ts
import { createHmac } from 'crypto';

const body = JSON.stringify({ room_id: "room01" });
const signature = createHmac("sha256", "SECRET")
    .update(body)
    .digest("hex");
```

###### Go

```go
const body = `{"room_id":"room01"}`
mac := hmac.New(sha256.New, []byte("SECRET"))
mac.Write([]byte(body))
signature := hex.EncodeToString(mac.Sum(nil))
```

###### Python

```python
import hmac
import hashlib

body = '{"room_id":"room01"}'
secret = b'SECRET'
signature = hmac.new(secret, body.encode(), hashlib.sha256).hexdigest()
```

###### Java

```java
import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

String body = "{\"room_id\":\"room01\"}";
String secret = "SECRET";
Mac sha256_HMAC = Mac.getInstance("HmacSHA256");
SecretKeySpec secret_key = new SecretKeySpec(secret.getBytes(), "HmacSHA256");
sha256_HMAC.init(secret_key);
String signature = javax.xml.bind.DatatypeConverter.printHexBinary(sha256_HMAC.doFinal(body.getBytes())).toLowerCase();
```

###### Ruby

```ruby
require 'openssl'

body = '{"room_id":"room01"}'
secret = 'SECRET'

digest = OpenSSL::Digest.new('sha256')
signature = OpenSSL::HMAC.hexdigest(digest, secret, body)
```

---

##### SDKs

For an even easier integration, you can use the following ready-to-use SDKs which handle the authentication process for you:

- [PHP](https://github.com/mynaparrot/plugNmeet-sdk-php)
- [JavaScript](https://github.com/mynaparrot/plugNmeet-sdk-js) ([NodeJS](https://www.npmjs.com/package/plugnmeet-sdk-js) & [Deno](https://deno.land/x/plugnmeet))
- [Python](https://github.com/vector-mj/plugnmeet-sdk) (Community Supported)

<a id="page-10"></a>

## حذف ضبط

حذف ضبط با شناسهٔ فایل انجام می‌شود. این عملیات باید پس از بررسی سیاست نگهداری، مجوز و وضعیت پردازش فایل انجام شود. رفتار پاسخ کاربردی از status و status_code بررسی و در تطبیق وضعیت محلی ثبت شود.

این بخش صفحهٔ `api/recording/delete` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/recording/delete) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/recording/delete.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/recording/delete.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Delete recording

Endpoint: `/recording/delete`

This endpoint allows you to permanently delete a recording from the server. This action is irreversible and will remove the recording file and all associated metadata.

This is a critical administrative function for managing storage and ensuring compliance with data retention policies. You would typically integrate this into your application's interface to allow authorized users to remove recordings they no longer need.

| Field     | Type   | Position | Required | Description                          |
| --------- | ------ | -------- | :------- | ------------------------------------ |
| record_id | string | root     | Yes      | Recording Id that you want to delete |

**Example**:

```json
{
  "record_id": "7f867cd7-7956-4a17-af46-6ddd4015a497-1761814595173"
}
```

##### Response

| Field  | Type    | Position | Description               |
| :----- |---------|----------| :------------------------ |
| status | boolean | root     | The status of the request |
| msg    | string  | root     | Response message          |
| status_code | string | root     | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |

<a id="page-11"></a>

## توکن دانلود ضبط

با record_id یک توکن موقت دانلود گرفته می‌شود و مسیر GET /download/recording/<TOKEN> فایل را تحویل می‌دهد. توکن برای مصرف فوری ساخته می‌شود و نباید به عنوان لینک دائمی در دیتابیس نگهداری شود. برای دانلود مجدد توکن تازه گرفته شود.

این بخش صفحهٔ `api/recording/download` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/recording/download) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/recording/download.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/recording/download.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Download Recording (Token)

Endpoint: `/recording/getDownloadToken`

This endpoint generates a secure, temporary token that can be used to download a recording file. Instead of providing a direct, permanent link to the recording, this API creates a short-lived, single-use token.

This is a security best practice that prevents unauthorized access and sharing of your recording files. The typical workflow is:
1. A user in your application clicks a "Download" button.
2. Your backend calls this API to get a token.
3. Your application constructs the download URL with the token and redirects the user to it.

The token will expire after a short period, ensuring that the download link cannot be reused or shared.

| Field     | Type   | Position | Required | Description                                 |
| --------- | ------ | -------- | :------- | ------------------------------------------- |
| record_id | string | root     | Yes      | The ID of the recording you want to download. |

**Example Request:**

```json
{
  "record_id": "7f867cd7-7956-4a17-af46-6ddd4015a497-1761814595173"
}
```

##### Response

| Field  | Type    | Position | Description                       |
| :----- | ------- | -------- | :------------------------------- |
| status | boolean | root     | Indicates if the request was successful. |
| msg    | string  | root     | Response message.                |
| status_code | string | root     | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |
| token  | string  | root     | Download token.                  |

Once you receive the token, you can construct the download URL using the following format:
```
https://Your-Plug-N-Meet-Server.com/download/recording/<TOKEN HERE>
```

<a id="page-12"></a>

## فهرست فایل‌های ضبط

ضبط‌ها با فیلترهای اتاق و نشست و صفحه‌بندی دریافت می‌شوند. اطلاعات فایل، metadata و زیرنویس در جدول‌های صفحه آمده است. یک نشست می‌تواند چند ضبط داشته باشد، بنابراین recording_id باید مستقل از room_id ذخیره شود.

این بخش صفحهٔ `api/recording/fetch` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/recording/fetch) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/recording/fetch.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/recording/fetch.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Fetch Recordings

Endpoint: `/recording/fetch`

This endpoint allows you to retrieve a paginated list of available recordings. Recordings are generated when a room session is recorded and are stored on the server after the session ends.

This API is essential for building a user-facing recording library where users can browse, view, and manage their past recordings. You can filter recordings by `room_id` or a specific room `sid` to narrow down the results, making it easy to display relevant content to your users.

##### Request Parameters

| Field    | Type   | Required | Description                                   |
| -------- | ------ | -------- | --------------------------------------------- |
| room_ids | array  | Yes      | Array of room IDs to fetch recordings for.    |
| room_sid   | string           | No       | Filter recordings from a specific room session by providing its `sid`.                                                                                                   |
| from     | number | No       | Starting index for records. Default is 0.     |
| limit    | number | No       | Maximum number of records to return. Default is 20. |
| order_by | string | No       | Sort order: `DESC` or `ASC`. Default is `DESC`. |

**Example Request:**

```json
{
  "room_ids": ["room01"],
  "from": 0,
  "limit": 20,
  "order_by": "DESC"
}
```

##### Response

| Field             | Type                        | Position | Description                                 |
| ----------------- | -------------------------- | -------- | ------------------------------------------- |
| status            | boolean                     | root     | Indicates if the request was successful.    |
| msg               | string                      | root     | Response message.                           |
| status_code | string | root     | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |
| [result](https://www.plugnmeet.org/docs/api/recording/fetch#result) | object                      | root     | Contains the recordings data.               |

###### Result

| Field              | Type                                      | Description                                 |
| ------------------ | ----------------------------------------- | ------------------------------------------- |
| total_recordings   | number                                    | Total number of recordings found.           |
| from               | number                                    | Starting index for the returned records.    |
| limit              | number                                    | Number of records returned.                 |
| order_by           | string                                    | Sort order used for the records.            |
| recordings_list    | Array\<[Recording Info](https://www.plugnmeet.org/docs/api/recording/fetch#recording-info)> | List of recording records.                  |

###### Recording Info

| Field              | Type   | Description                                 |
| ------------------ | ------ | ------------------------------------------- |
| record_id          | string | Unique identifier for the recording.        |
| room_id            | string | ID of the room associated with the recording.|
| room_sid           | string | SID of the room.                            |
| file_path          | string | Path to the recording file.                 |
| file_size          | number | Size of the recording file in bytes.        |
| creation_time      | number | Recording creation time (Unix timestamp).   |
| room_creation_time | number | Room creation time (Unix timestamp).        |
| metadata           | [RecordingMetadata](https://www.plugnmeet.org/docs/api/recording/fetch#recordingmetadata-object) | Metadata associated with the recording.     |

###### RecordingMetadata Object

| Field       | Type                                               | Description                                                                |
| ----------- |----------------------------------------------------|----------------------------------------------------------------------------|
| title       | string                                             | The title of the recording.                                                |
| description | string                                             | The description of the recording.                                          |
| subtitles   | map [RecordingSubtitle](https://www.plugnmeet.org/docs/api/recording/fetch#recordingsubtitle-object) | A map of language codes to subtitle objects. e.g. `{"en": {"url": "..."}}` |
| extra_data  | map                               | A map of extra data stored with the recording. e.g. `{"key": "value"}`     |

###### RecordingSubtitle Object

| Field | Type   | Description                  |
| ----- | ------ | ---------------------------- |
| label | string | The label for the subtitle.  |
| url   | string | The URL of the subtitle file.|

<a id="page-13"></a>

## ادغام ضبط‌های چندبخشی

ادغام به صورت غیرهم‌زمان اجرا می‌شود. انتخاب فایل‌ها با by_session یا by_ids انجام می‌شود و خروجی جدید پس از پردازش از طریق وب‌هوک آماده‌شدن ضبط پیگیری می‌شود. شناسهٔ ضبط ادغام‌شده و زمان‌بندی آن باید جداگانه ثبت شود تا بازهٔ گفتار و قطعه‌های اصلی قابل ردیابی بمانند.

این بخش صفحهٔ `api/recording/merge-recordings` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/recording/merge-recordings) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/recording/merge-recordings.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/recording/merge-recordings.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Merge Recordings

Endpoint: `/recording/mergeRecordings`

In long-running sessions, it's common for administrators or presenters to stop and restart a recording, resulting in multiple separate video files. This API provides a powerful solution to consolidate these segments into one continuous recording, improving the viewing experience.

Using FFmpeg's append feature, this endpoint stitches together recording files into a new recording with a unique `record_id`. The original recording files are not deleted, giving you the flexibility to manage or discard them as needed. This is ideal for creating a polished, final version of a session that was recorded in parts.

##### Asynchronous Operation & Webhook Notification

This API operation is asynchronous. When you call this endpoint, the request is immediately assigned to a `plugnmeet-recorder` instance to be processed in the background. This requires you to have at least one active recorder running in either `both` or `transcoderOnly` mode. The API will instantly return a confirmation that the job has been scheduled.

Because the merging happens in the background, the final result is delivered via a webhook. Once the new recording file is ready, your server will receive a `recording_proceeded` event containing the `record_id` and other metadata for the newly merged file, just like a standard recording.

##### Use Cases
*   **Creating a single file for long webinars:** If you pause a recording during a break, you can merge the parts for a seamless replay.
*   **Combining training session modules:** Merge separate recorded modules into one comprehensive training video.
*   **Post-production simplification:** Provide a single file for easier editing and distribution.

##### Request Parameters

The request body must contain one of the following scopes: `by_session` or `by_ids`.

###### `by_session`

Merge all recordings from a single room session. The recordings are appended chronologically in ascending order (ASC) as they were saved, ensuring a natural sequence.

| Field                 | Type           | Required | Description                                                                                                                                                             |
| --------------------- | -------------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| room_sid              | string         | Yes      | The unique session identifier (`sid`) of the room whose recordings you want to merge. All recordings under this `room_sid` will be considered for the merge.              |
| exclude_recording_ids | array (string) | No       | An array of `record_id`s to exclude from the merge. This is useful if you need to omit specific segments, such as a false start or an unwanted portion of the session. |

###### `by_ids`

Merge a specific list of recordings. The final video will be generated by appending recordings in the order they are provided in the `recording_ids` array.

| Field         | Type           | Required | Description                                                                                                                                                           |
| ------------- | -------------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| room_id       | string         | Yes      | The `room_id` that all `recording_ids` belong to.                                                                                                                     |
| recording_ids | array (string) | Yes      | A list of `record_id`s to merge. All recordings must belong to the specified `room_id`, but can be from different sessions. The list must contain at least two `record_id`s. |

##### Example Requests

###### Example 1: Merge by Session

This example merges all recordings from the session `c5b79ab5...` except for one excluded recording.

```json
{
  "by_session": {
    "room_sid": "c5b79ab5-b832-4972-9f9b-ba1f1e3369bd",
    "exclude_recording_ids": [
      "c5b79ab5-b832-4972-9f9b-ba1f1e3369bd-1779521611869"
    ]
  }
}
```

###### Example 2: Merge by IDs

This example merges two specific recordings from the room `room01` in the specified order.

```json
{
  "by_ids": {
    "room_id": "room01",
    "recording_ids": [
      "c5b79ab5-b832-4972-9f9b-ba1f1e3369bd-1779521611869",
      "c5b79ab5-b832-4972-9f9b-ba1f1e3369bd-1779521613869"
    ]
  }
}
```

##### Response

The API provides an immediate response indicating that the merge job has been successfully scheduled. The final result of the merge operation will be delivered via a `recording_proceeded` webhook event once the process is complete.

| Field       | Type    | Description                                                                                                                              |
| ----------- | ------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| status      | boolean | Indicates if the request was successful.                                                                                                 |
| msg         | string  | A message confirming the job was queued.                                                                                                 |
| status_code | number  | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |

<a id="page-14"></a>

## جزئیات یک فایل ضبط

با record_id مشخصات ضبط دریافت می‌شود. Core باید ضبط را به نشست و ارائه‌دهندهٔ صحیح متصل کند و سپس قابلیت پخش یا دانلود را مطابق سطح دسترسی همان جلسه نمایش دهد.

این بخش صفحهٔ `api/recording/recording-info` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/recording/recording-info) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/recording/recording-info.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/recording/recording-info.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Recording info

Endpoint: `/recording/info`

This endpoint retrieves all the metadata and details associated with a specific recording. By providing a `record_id`, you can get comprehensive information about the recording, including its file size, creation time, and associated room details.

This is useful for displaying detailed information about a recording in your application's interface, such as on a dedicated recording playback page.

| Field     | Type   | Position | Required | Description                          |
| --------- | ------ | -------- | :------- | ------------------------------------ |
| record_id | string | root     | Yes      | Recording Id that you want to get details |

**Example**:

```json
{
  "record_id": "7f867cd7-7956-4a17-af46-6ddd4015a497-1761814595173"
}
```

##### Response

| Field                                                          | Type    | Position | Description               |
|:---------------------------------------------------------------| ------- | -------- | :------------------------ |
| status                                                         | boolean | root     | The status of the request |
| msg                                                            | string  | root     | Response message          |
| status_code | string | root     | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |
| [recording_info](https://www.plugnmeet.org/docs/api/recording/fetch#recording-info)     | obect\<[Recording Info](https://www.plugnmeet.org/docs/api/recording/fetch#recording-info)> | root     | Response message          |
| [room_info](https://www.plugnmeet.org/docs/api/room/fetch-past#past-room-info)          | obect\<[past-room-info](https://www.plugnmeet.org/docs/api/room/fetch-past#past-room-info)> | root     | Response message          |

<a id="page-15"></a>

## به‌روزرسانی metadata ضبط

این متد metadata فایل ضبط را اصلاح می‌کند. نام فیلدها و ساختار مورد قبول از جدول پیوست گرفته شود. دادهٔ محلی آی‌نشست و metadata سمت ارائه‌دهنده باید بر اساس مالکیت هر فیلد به‌روزرسانی شوند.

این بخش صفحهٔ `api/recording/update-metadata` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/recording/update-metadata) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/recording/update-metadata.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/recording/update-metadata.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Update Recording Metadata

Endpoint: `/recording/updateMetadata`

This API allows you to update the metadata of a specific recording. It intelligently handles partial updates based on the provided fields:
- To update a field, provide a new value.
- To clear a text field (like Title), provide an empty string "".
- To clear a specific map entry (like a subtitle), provide an empty object for that key.
- If a field is omitted (i.e., nil), its existing value is kept.

| Field     | Type   | Position | Required | Description                          |
| --------- | ------ | -------- | :------- | ------------------------------------ |
| record_id | string | root     | Yes      | The ID of the recording to update. |
| metadata  | [RecordingMetadata](https://www.plugnmeet.org/docs/api/recording/fetch#recordingmetadata-object) | root     | No       | The metadata to update.   |

**Example Request:**

```json
{
  "record_id": "7f867cd7-7956-4a17-af46-6ddd4015a497-1761814595173",
  "metadata": {
    "title": "New Recording Title",
    "description": "Intro session",
    "subtitles": {
      "en": {
        "label": "English",
        "url": "https://example.com/subtitle.vtt"
      }
    },
    "extra_data": {
      "any_key": "any_value",
      "source": "cloud"
    }
  }
}
```

##### Response

| Field  | Type    | Position | Description               |
| :----- | ------- | -------- | :------------------------ |
| status | boolean | root     | The status of the request |
| msg    | string  | root     | Response message          |
| status_code | string | root     | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |

<a id="page-16"></a>

## ارسال پیام و اعلان به اتاق

این API پیام گفتگو یا اعلان را به مخاطبان انتخابی می‌فرستد. هدف‌گیری عمومی، کاربران مشخص و مدیران طبق قواعد جدول صفحه تعیین می‌شود. ساختار ChatMessage و NotificationMsg جداگانه است و مثال‌های ارسال عمومی و خصوصی در پیوست آمده‌اند.

این بخش صفحهٔ `api/room/broadcast-to-room` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/room/broadcast-to-room) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/room/broadcast-to-room.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/room/broadcast-to-room.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Broadcast to Room

Endpoint: `/room/broadcastToRoom`

This powerful API allows your backend server to inject messages or notifications directly into an active Plug-N-Meet session in real-time. You can use it to send system-wide announcements, deliver private messages to specific users, or trigger custom notifications from your application's logic.

This endpoint is ideal for building integrations that require server-to-client communication, such as:
*   Sending alerts based on external events.
*   Delivering automated messages or instructions to participants.
*   Creating custom, private notification systems within your application.

##### Request Parameters

The request must contain either a `chat_msg` or a `notification_msg` object.

| Field            | Type    | Required | Description                                                                                             |
| ---------------- | ------- | -------- | ------------------------------------------------------------------------------------------------------- |
| room_id          | string  | Yes      | The unique identifier of the active room to which you want to broadcast.                                |
| only_to_admins   | boolean | No       | If `true`, the message will be sent to all participants with admin/moderator privileges.                |
| to_user_id       | string  | No       | If provided, the message will be sent to the specified user.                                            |
| chat_msg         | object  | No       | A [chat message](https://www.plugnmeet.org/docs/api/room/broadcast-to-room#chatmessage) object to be sent.                                                       |
| notification_msg | object  | No       | A [notification message](https://www.plugnmeet.org/docs/api/room/broadcast-to-room#notificationmsg) object to be displayed.                                      |

###### Targeting Logic

The delivery of the message is determined by the `only_to_admins` and `to_user_id` fields:
-   If **neither** field is set, the message is sent to **everyone** in the room.
-   If `only_to_admins` is `true`, the message is sent to **all admins**.
-   If `to_user_id` is provided, the message is sent to that **specific user**.
-   If **both** `only_to_admins` and `to_user_id` are provided, the message is sent to **all admins AND the specified user**.

###### ChatMessage

This object represents a message that will appear in the room's chat panel.

| Field   | Type   | Required | Description                      |
| ------- | ------ | -------- | -------------------------------- |
| message | string | Yes      | The content of the chat message. |

###### NotificationMsg

This object represents a system-level notification that will pop up on the user's screen.

| Field      | Type                                                                                                                           | Required | Description                                                                 |
| ---------- | ------------------------------------------------------------------------------------------------------------------------------ | -------- | --------------------------------------------------------------------------- |
| text       | string                                                                                                                         | Yes      | The text content of the notification.                                       |
| type       | [NatsSystemNotificationTypes](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_nats_msg.proto#L142) | No       | The style of the notification. `0 = info`; `1 = warning`; `3 = error`       |
| with_sound | boolean                                                                                                                        | No       | If `true`, the notification will play a sound for the recipient.            |

##### Example

###### Example 1: Sending a Public Chat Message to Everyone

```json
{
  "room_id": "room01",
  "chat_msg": {
    "message": "Hello everyone, the webinar will begin in 5 minutes."
  }
}
```

###### Example 2: Sending a Private Notification to a Specific User

```json
{
  "room_id": "room01",
  "to_user_id": "user-456",
  "notification_msg": {
    "text": "Your private document is ready for review.",
    "type": 0
  }
}
```

###### Example 3: Sending a Notification to All Admins and One Specific User

```json
{
  "room_id": "room01",
  "only_to_admins": true,
  "to_user_id": "vip-user-789",
  "notification_msg": {
    "text": "A VIP has just joined the waiting room.",
    "type": 1,
    "with_sound": true
  }
}
```

##### Response

| Field       | Type    | Description                              |
| ----------- | ------- | ---------------------------------------- |
| status      | boolean | Indicates if the request was successful. |
| msg         | string  | Response message.                        |
| status_code | string  | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |

<a id="page-17"></a>

## ساخت نظرسنجی و آزمون

نظرسنجی در اتاق فعال با گزینه‌ها و تنظیمات انتخاب چندگانه، ناشناس یا آزمون زمان‌دار ساخته می‌شود. جدول درخواست، ساختار گزینه و مثال‌های سه حالت در پیوست آمده‌اند. خطاهای اعتبارسنجی تعریف‌شدهٔ همین صفحه در اتصال بیرونی بررسی شوند.

این بخش صفحهٔ `api/room/create-poll` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/room/create-poll) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/room/create-poll.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/room/create-poll.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Create Poll

Endpoint: `/room/createPoll`

This API allows your backend server to push a complete poll into an active Plug-N-Meet session in real time. While moderators can create polls from within the client, this endpoint lets an external application inject a poll directly into a running session, making it a powerful building block for automated integrations.

This endpoint is ideal for building integrations such as:
*   Launching a live quiz in a classroom from an external learning platform.
*   Collecting audience votes triggered by events in your application.
*   Scheduling recurring polls or surveys from an external tool.

Once the poll is pushed, it is created in the running state immediately:
1.  Every participant in the live room receives the new-poll notification and can vote right away.
2.  The moderator can then close it, publish the results, or reopen it — just like any poll created from within the client.
3.  If `duration` is set, the poll closes automatically when the time runs out.

For this API call to succeed, the session (`room_id`) must be a currently active room. The request fails if the room has not been created or has already ended.

##### Request Parameters

| Field        | Type    | Required | Description                                                                                                                              |
| ------------ | ------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| room_id      | string  | Yes      | The unique identifier of the active room into which you want to push the poll.                                                            |
| user_id      | string  | No       | The user ID to attribute the poll to. If omitted, it defaults to `external-api`.                                                          |
| question     | string  | Yes      | The poll question. Must not be empty.                                                                                                     |
| [options](https://www.plugnmeet.org/docs/api/room/create-poll#poll-option) | array | Yes | The poll options. At least two are required and each must have non-empty text. See [Poll Option](https://www.plugnmeet.org/docs/api/room/create-poll#poll-option).                           |
| is_anonymous | boolean | No       | If `true`, voting is anonymous: individual voter choices are never stored per-user. Default: `false`.                                     |
| is_multiple  | boolean | No       | If `true`, participants can select multiple options when voting. Default: `false`.                                                        |
| is_quiz      | boolean | No       | If `true`, the poll runs as a quiz. Correct answers are hidden while the quiz is running and revealed in the results after it is closed. A quiz requires at least one option marked as correct. Default: `false`. |
| duration     | number  | No       | Auto-close duration in seconds. If greater than `0`, the poll closes automatically when the time runs out. Maximum: `3600` (60 minutes). Default: `0` (no time limit). |

###### Poll Option

Each entry in the `options` array represents one selectable choice.

| Field      | Type    | Required | Description                                                                    |
| ---------- | ------- | -------- | ------------------------------------------------------------------------------ |
| id         | number  | Yes      | A simple sequential number identifying the option (1, 2, 3, ...).              |
| text       | string  | Yes      | The option text. Must not be empty.                                            |
| is_correct | boolean | No       | If `true`, this option is a correct answer. Only meaningful when `is_quiz` is `true`. |

##### Example

###### Example 1: A Simple Poll

```json
{
  "room_id": "room01",
  "question": "Which topic should we cover next?",
  "options": [
    {
      "id": 1,
      "text": "Advanced whiteboard features"
    },
    {
      "id": 2,
      "text": "Recording and playback"
    }
  ]
}
```

###### Example 2: An Anonymous Multi-Select Poll

```json
{
  "room_id": "room01",
  "user_id": "user-42",
  "question": "Which sessions did you attend today?",
  "options": [
    {
      "id": 1,
      "text": "Morning keynote"
    },
    {
      "id": 2,
      "text": "Product workshop"
    },
    {
      "id": 3,
      "text": "Networking session"
    }
  ],
  "is_anonymous": true,
  "is_multiple": true
}
```

###### Example 3: A Timed Quiz

```json
{
  "room_id": "room01",
  "question": "What does the HASH-SIGNATURE header contain?",
  "options": [
    {
      "id": 1,
      "text": "An HMAC-SHA256 signature of the request body",
      "is_correct": true
    },
    {
      "id": 2,
      "text": "The API key"
    }
  ],
  "is_quiz": true,
  "duration": 120
}
```

##### Response

| Field       | Type    | Description                              |
| ----------- | ------- | ---------------------------------------- |
| status      | boolean | Indicates if the request was successful. |
| msg         | string  | Response message.                        |
| poll_id     | string  | The unique ID of the newly created poll. |
| status_code | string  | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |

##### Error Responses

On failure, `status` is `false` and `msg` describes the problem. Validation failures return stable, machine-readable keys that the client UI can translate directly:

| msg key                           | Description                                                              |
| --------------------------------- | ------------------------------------------------------------------------ |
| `polls.errors.question-required`  | The `question` field is missing or empty.                                |
| `polls.errors.min-options`        | Fewer than two options were provided.                                    |
| `polls.errors.option-required`    | An option is missing its text.                                           |
| `polls.errors.quiz-needs-correct` | The poll is marked as a quiz, but no option is marked as correct.        |
| `polls.errors.duration-cap`       | The `duration` exceeds the maximum of 3600 seconds (60 minutes).         |

Other common failures are returned as plain messages: a missing `room_id`, or a room that is not currently active (this endpoint requires the session to be running at the moment the poll is pushed).

<a id="page-18"></a>

## ساخت اتاق و تنظیم قابلیت‌ها

اتاق با room_id پایدار و metadata ساخته می‌شود. تنظیمات room_features، ضبط، گفتگو، وایت‌برد، انتظار، اتاق فرعی، نظرسنجی، AI، SIP، E2EE و قفل‌های پیش‌فرض به صورت بخش‌های تو در تو هستند. برای آی‌نشست analytics به‌طور صریح فعال و نقش‌ها در join کنترل شوند. گزینهٔ ضبط خودکار و اجازهٔ ضبط باید مطابق سیاست همان جلسه تعیین شود. همهٔ جدول‌های تو در تو و مثال کامل درخواست در پیوست این بخش حفظ شده‌اند.

این بخش صفحهٔ `api/room/create` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/room/create) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/room/create.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/room/create.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Create Room

Endpoint: `/room/create`

Before you create your first room, it's helpful to understand how rooms work in Plug-N-Meet. Think of a room not as a permanent space, but as a **temporary live session**.

Here's the typical lifecycle:
1.  You create a room using this API endpoint.
2.  You generate access tokens for users to join the session.
3.  The session remains active as long as participants are present.
4.  The session ends automatically when the last participant leaves or can be terminated via an API call.

Once a session is over, the room is finalized, and all associated data (like chat messages and user lists) is cleared. This ensures that each new session starts fresh.

> **Pro Tip: Simulating Permanent Rooms**
>
> If you want to create the experience of a "permanent" room that users can join at any time, you can build this logic into your application.
>
> When a user attempts to join, your application should:
> 1.  Check if an active session for that `room_id` already exists (e.g., using the `isRoomActive` API).
> 2.  If no session exists, call this `create` endpoint to start a new one.
> 3.  Finally, generate an access token to allow the user to join.
>
> This approach gives you the flexibility of persistent-like rooms while leveraging Plug-N-Meet's temporary session model.

##### Request Parameters

| Field                 | Type   | Required | Description                                                                            |
| --------------------- | ------ | -------- |----------------------------------------------------------------------------------------|
| room_id               | string | Yes      | A unique identifier for the room. Since rooms are temporary, you can reuse a `room_id` after a session has ended. |
| max_participants      | number | No       | The maximum number of participants allowed in the room.                                |
| empty_timeout         | number | No       | The number of seconds the room will remain active after creation if no one joins.      |
| [metadata](https://www.plugnmeet.org/docs/api/room/create#metadata) | object | Yes      | Additional room configuration details.                                                 |

###### Metadata

| Field                                           | Type                | Required | Description                                                                                      |
| ----------------------------------------------- |---------------------| -------- | ------------------------------------------------------------------------------------------------ |
| room_title                                      | string              | Yes      | The title of the room or meeting.                                                                |
| welcome_message                                 | string              | No       | A message displayed to participants when they join.                                              |
| webhook_url                                     | string              | No       | URL to receive webhook events from Plug-N-Meet.                                                  |
| logout_url                                      | string              | No       | URL to redirect users after the meeting or session ends.                                         |
| [room_features](https://www.plugnmeet.org/docs/api/room/create#room-features)                 | object              | Yes      | Settings to enable or disable various room features.                                             |
| [default_lock_settings](https://www.plugnmeet.org/docs/api/room/create#default-lock-settings) | object              | No       | Default settings to lock specific features for users.                                            |
| [copyright_conf](https://www.plugnmeet.org/docs/api/room/create#copyright-configuration)             | object              | No       | Copyright configuration.                                                                         |
| extra_data                                      | map | No       | A map of custom key-value pairs. Both keys and values must be strings. e.g. `{"key": "value"}`|.                                                           |

###### Room Features

| Field                                                                | Type    | Required | Description                                                                                            |
| -------------------------------------------------------------------- | ------- | -------- | ------------------------------------------------------------------------------------------------------ |
| allow_webcams                                                        | boolean | Yes      | Enable or disable webcam support.                                                                      |
| mute_on_start                                                        | boolean | Yes      | Automatically mute microphones when participants join.                                                 |
| allow_screen_share                                                   | boolean | Yes      | Enable or disable screen sharing.                                                                      |
| admin_only_webcams                                                   | boolean | Yes      | Allow webcam access only for admins.                                                                   |
| allow_view_other_webcams                                             | boolean | Yes      | Allow participants to view each other's webcams. If `false`, only moderators can see all webcams. |
| allow_view_other_users_list                                          | boolean | Yes      | Restrict viewing of the user list to moderators only.                                                  |
| enable_analytics                                                     | boolean | No       | Enable or disable analytics reporting for the session. Default: false                                  |
| allow_virtual_bg                                                     | boolean | No       | Enable or disable virtual background options.                                                          |
| allow_raise_hand                                                     | boolean | No       | Enable or disable the "raise hand" feature.                                                            |
| auto_gen_user_id                                                     | boolean | No       | PlugNmeet requires a unique `userId` for each participant. Enable this if you prefer not to manage `userId`s or need to allow the same user to join from multiple devices. When `true`, PlugNmeet generates a unique `user_id` for each session. Any `user_id` you provide will be stored as `ex_user_id`, retrievable via the `getActiveRoomInfo` API. Default: `false` |
| room_duration                                                        | number  | No       | Set a fixed duration for the room in minutes. 0 means unlimited.                                       |
| [recording_features](https://www.plugnmeet.org/docs/api/room/create#recording-features)                            | object  | Yes      | Recording settings.                                                                                    |
| [chat_features](https://www.plugnmeet.org/docs/api/room/create#chat-features)                                      | object  | Yes      | Chat settings.                                                                                         |
| [shared_note_pad_features](https://www.plugnmeet.org/docs/api/room/create#shared-notepad-features)                | object  | Yes      | Shared notepad settings.                                                                               |
| [whiteboard_features](https://www.plugnmeet.org/docs/api/room/create#whiteboard-features)                          | object  | Yes      | Whiteboard settings.                                                                                   |
| [external_media_player_features](https://www.plugnmeet.org/docs/api/room/create#external-media-player-features)    | object  | Yes      | External media player settings.                                                                        |
| [external_broadcasting_features](https://www.plugnmeet.org/docs/api/room/create#external-broadcasting-features)    | object  | No       | External broadcasting settings.                                                                        |
| [waiting_room_features](https://www.plugnmeet.org/docs/api/room/create#waiting-room-features)                      | object  | Yes      | Waiting room settings.                                                                                 |
| [breakout_room_features](https://www.plugnmeet.org/docs/api/room/create#breakout-room-features)                    | object  | Yes      | Breakout room settings.                                                                                |
| [display_external_link_features](https://www.plugnmeet.org/docs/api/room/create#display-external-link-features)    | object  | Yes      | Settings for displaying external links.                                                                |
| [ingress_features](https://www.plugnmeet.org/docs/api/room/create#ingress-features)                                | object  | No       | RTMP ingress settings.                                                                                 |
| [polls_features](https://www.plugnmeet.org/docs/api/room/create#polls-features)                                    | object  | No       | Polls settings.                                                                                        |
| [insights_features](https://www.plugnmeet.org/docs/api/room/create#insights-features)                              | object  | No       | AI-powered insights, transcription, and translation settings.                                          |
| [sip_dial_in_features](https://www.plugnmeet.org/docs/api/room/create#sip-dial-in-features)                        | object  | No       | SIP dial-in settings.                                                                                  |
| [end_to_end_encryption_features](https://www.plugnmeet.org/docs/api/room/create#end-to-end-encryption-e2ee-features) | object  | No       | End-to-End Encryption (E2EE) settings.                                                                 |

###### Recording Features
<details>
<summary>View Details</summary>

| Field                       | Type    | Required | Description                                                                           |
| --------------------------- | ------- | -------- | ------------------------------------------------------------------------------------- |
| is_allow                    | boolean | Yes      | Enable or disable recording for the meeting.                                          |
| is_allow_cloud              | boolean | Yes      | Enable or disable cloud recording.                                                    |
| is_allow_local              | boolean | Yes      | Enable or disable local recording.                                                    |
| enable_auto_cloud_recording | boolean | No       | Automatically start cloud recording when a moderator or admin joins.                  |
| only_record_admin_webcams   | boolean | No       | If `true`, only the webcams of admin users will be included in the recording.           |
| [recorder_bot_options](https://www.plugnmeet.org/docs/api/room/create#recorder-bot-options) | object | No | Advanced options for the recorder bot.                                        |

</details>

###### Chat Features
<details>
<summary>View Details</summary>

| Field               | Type    | Required | Description                             |
| ------------------- | ------- | -------- | --------------------------------------- |
| is_allow            | boolean | Yes      | Enable or disable chat for the meeting. |
| is_allow_file_upload | boolean | Yes      | Enable or disable file uploads in chat. |

</details>

###### Shared Notepad Features
<details>
<summary>View Details</summary>

| Field    | Type    | Required | Description                                   |
| -------- | ------- | -------- | --------------------------------------------- |
| is_allow | boolean | Yes      | Enable or disable the shared notepad feature. |

</details>

###### Whiteboard Features
<details>
<summary>View Details</summary>

| Field        | Type    | Required | Description                                                                                                   |
| ------------ | ------- | -------- | ------------------------------------------------------------------------------------------------------------- |
| is_allow     | boolean | Yes      | Enable or disable the whiteboard feature.                                                                     |
| preload_file | string  | No       | Preload a remote (http/https) presentation file for the whiteboard. The file must be directly accessible without redirection. |

</details>

###### External Media Player Features
<details>
<summary>View Details</summary>

| Field    | Type    | Required | Description                                                                              |
| -------- | ------- | -------- | ---------------------------------------------------------------------------------------- |
| is_allow | boolean | Yes      | Enable or disable playback of video/audio from external sources. Moderators can also upload local media. |

</details>

###### External Broadcasting Features
<details>
<summary>View Details</summary>

| Field                | Type   | Required | Description                                                              |
| -------------------- | ------ | -------- | ------------------------------------------------------------------------ |
| is_allow             | boolean| Yes      | Enable or disable external broadcasting features.                        |
| is_allow_rtmp        | boolean| Yes      | When enabled, moderators can start live streaming the room to external platforms like YouTube. |
| [recorder_bot_options](https://www.plugnmeet.org/docs/api/room/create#recorder-bot-options) | object | No | Advanced options for the recorder bot used for broadcasting. |

</details>

###### Waiting Room Features
<details>
<summary>View Details</summary>

| Field     | Type    | Required | Description                                                                                                     |
| --------- | ------- | -------- | --------------------------------------------------------------------------------------------------------------- |
| is_active | boolean | Yes      | Enable to activate the waiting room feature. Users remain in the waiting room until allowed to join by a moderator. |

</details>

###### Breakout Room Features
<details>
<summary>View Details</summary>

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| is_allow | boolean | Yes | Enable or disable breakout rooms. |
| allowed_number_rooms | number | No | Maximum number of breakout rooms that can be created simultaneously. Default: 6. |
| allow_return_to_main_room | boolean | No | Allow participants to return to the main room from a breakout room. Default: `true`. |
| allow_self_select | boolean | No | Allow participants to self-select which breakout room to join. Default: `false`. |
| preassigned_rooms | array | No | Pre-assign participants to specific breakout rooms. |

###### Preassigned Breakout Room

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| title | string | Yes | The title of the pre-assigned breakout room. |
| user_ids | array(string) | No | A list of user IDs to be automatically assigned to this breakout room. |

</details>

###### Display External Link Features
<details>
<summary>View Details</summary>

| Field    | Type    | Required | Description                                                                                                                                                                                                                                                                                       |
| -------- | ------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| is_allow | boolean | Yes      | Enable or disable the display of external links inside an iframe. Useful for showing websites, quizzes, or games. **Note:** The external site must allow being embedded in an iframe. Moderators can pass values like `name`, `userId`, `role`, and `meetingId` to the URL. |

</details>

###### Ingress Features
<details>
<summary>View Details</summary>

| Field    | Type    | Required | Description                                                                                                   |
| -------- | ------- | -------- | ------------------------------------------------------------------------------------------------------------- |
| is_allow | boolean | Yes      | Enable media ingress to stream content directly into the session. This is useful for professional broadcasting using software like OBS Studio. Plug-N-Meet supports both RTMP and WHIP (WebRTC-HTTP Ingestion Protocol) for low-latency streaming. |

</details>

###### Polls Features
<details>
<summary>View Details</summary>

| Field    | Type    | Required | Description                       |
| -------- | ------- | -------- | --------------------------------- |
| is_allow | boolean | Yes      | Enable or disable polls features. |

</details>

###### Insights Features
<details>
<summary>View Details</summary>

**Server Configuration Required**
For any of these features to work, you must first configure the `insights` section of your server's main `config.yaml` file.

The configuration has a powerful two-part structure:
1.  First, you define your **provider accounts**, such as `azure` or `google`, and add your API credentials for each.
2.  Second, you map each **service** (e.g., `transcription`, `ai_text_chat`) to a specific provider account.

This flexible design allows you to use different providers for different tasks (e.g., Azure for translation and Google for AI chat) and even select different models for cost-performance optimization (e.g., a powerful model for chat and a cheaper one for summarization). Please refer to the server configuration guide for detailed setup instructions.


This object controls all AI-powered features provided by the Insights Platform.

| Field                                       | Type   | Required | Description                                           |
| ------------------------------------------- | ------ | -------- | ----------------------------------------------------- |
| is_allow                                    | boolean| Yes      | Enable or disable all Insights Platform features.     |
| [transcription_features](https://www.plugnmeet.org/docs/api/room/create#transcription-features) | object | No       | Settings for speech-to-text and spoken translation.   |
| [chat_translation_features](https://www.plugnmeet.org/docs/api/room/create#chat-translation-features) | object | No       | Settings for chat translation.                        |
| [ai_features](https://www.plugnmeet.org/docs/api/room/create#ai-features)                 | object | No       | Settings for generative AI features.                  |

###### Transcription Features

| Field                       | Type    | Required | Description                                                 |
| --------------------------- | ------- | -------- | ----------------------------------------------------------- |
| is_allow                    | boolean | Yes      | Enable or disable real-time transcription.                  |
| is_allow_translation        | boolean | No       | Enable or disable translation of the transcription.         |
| is_allow_speech_synthesis   | boolean | No       | Enable or disable spoken translations (text-to-speech).     |

###### Chat Translation Features

| Field    | Type    | Required | Description                                           |
| -------- | ------- | -------- | ----------------------------------------------------- |
| is_allow | boolean | Yes      | Enable or disable real-time translation of chat messages. |

###### AI Features

| Field                                              | Type   | Required | Description                                |
| -------------------------------------------------- | ------ | -------- | ------------------------------------------ |
| is_allow                                           | boolean| Yes      | Enable or disable all generative AI features.  |
| [ai_text_chat_features](https://www.plugnmeet.org/docs/api/room/create#ai-text-chat-features)    | object | No       | Settings for the AI Chat Assistant.        |
| [meeting_summarization_features](https://www.plugnmeet.org/docs/api/room/create#meeting-summarization-features) | object | No       | Settings for AI-powered meeting summaries. |

###### AI Text Chat Features

| Field    | Type    | Required | Description                          |
| -------- | ------- | -------- | ------------------------------------ |
| is_allow | boolean | Yes      | Enable or disable the AI Chat Assistant. |

###### Meeting Summarization Features

| Field    | Type    | Required | Description                                      |
| -------- | ------- | -------- | ------------------------------------------------ |
| is_allow | boolean | Yes      | Enable or disable the generation of meeting summaries. |

</details>

###### SIP Dial-In Features
<details>
<summary>View Details</summary>

**Server Configuration Required**
For SIP dial-in to function, you must first configure the `livekit_sip_info` section of your server's main `config.yaml` file. This involves setting up the SIP gateway (e.g., `livekit/sip`) and integrating it with an external SIP trunking provider.

Please refer to the server configuration guide and the [livekit/sip GitHub repository](https://github.com/livekit/sip) for detailed setup instructions.


| Field                      | Type    | Required | Description                                                                                                                                                                                                                                                                                                                                                                                                                    |
| -------------------------- | ------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| is_allow                   | boolean | Yes      | Enable or disable SIP dial-in for the room.                                                                                                                                                                                                                                                                                                                                                                                    |
| enable_dial_in_on_create   | boolean | No       | If `true`, the SIP dial-in service will automatically start when the room is created. If `false`, a moderator must manually start it from within the meeting. Default: `false`.                                                                                                                                                                                                                                                |
| hide_phone_number          | boolean | No       | If `true`, the dial-in phone number will be masked in the participant list, showing only the last 4 digits. Default: `false`.                                                                                                                                                                                                                                                                                                   |

</details>

###### End-to-End Encryption (E2EE) Features
<details>
<summary>View Details</summary>

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| is_enabled | boolean | Yes | Enable or disable E2EE. Supported browsers: Chromium 83+, Google Chrome, Microsoft Edge, Safari, Firefox 117+. **Note:** Users cannot join if their browser does not support E2EE. |
| included_chat_messages | boolean | No | Enable or disable E2EE for chat messages. |
| included_whiteboard | boolean | No | Enable or disable E2EE for whiteboard messages (SCENE_UPDATE, POINTER_UPDATE). May increase CPU usage; enable only if necessary. |
| encryption_key | string | No | Optional. If you provide a key, it will be used as the E2EE key. If not provided, a key will be generated automatically. The key must be a string of any length. |
| enabled_self_insert_encryption_key | boolean | No | If `true`, users will be prompted to enter a secret key upon joining. This enables a zero-trust model where the key is **never sent to the server**. If `false`, the server generates and distributes the key. Default: `false`. <br/><br/> **Important:** When this is `true`, all audio-based AI features (like transcription and summarization) are automatically disabled. This is because the server does not have the key and therefore cannot provide it to the AI agent to access the audio stream. |

</details>

###### Recorder Bot Options
<details>
<summary>View Details</summary>

| Field                         | Type   | Required | Description                                                              |
| ----------------------------- | ------ | -------- | ------------------------------------------------------------------------ |
| enable_auto_close_chat_panel  | boolean| No       | If `true`, the chat panel will be closed automatically after a certain duration. |
| duration_after_last_message   | number | No       | Duration in seconds. After this duration, the chat panel will be closed. |

</details>

###### Default Lock Settings
<details>
<summary>View Details</summary>

| Field                  | Type    | Required | Description                  |
| ---------------------- | ------- | -------- | ---------------------------- |
| lock_microphone        | boolean | No       | Lock microphone for users.   |
| lock_webcam            | boolean | No       | Lock webcam for users.       |
| lock_screen_sharing    | boolean | No       | Lock screen sharing for users.|
| lock_chat              | boolean | No       | Lock chat for users.         |
| lock_chat_send_message | boolean | No       | Lock sending messages in chat.|
| lock_chat_file_share   | boolean | No       | Lock file sharing in chat.   |
| lock_private_chat      | boolean | No       | Lock private chat for users. |

</details>

###### Copyright Configuration
<details>
<summary>View Details</summary>

This feature is available only if the server configuration `client > copyright_conf > allow_override` is set to `true`.

| Field   | Type   | Required | Description                                                                 |
| ------- | ------ | -------- | ----------------------------------------------------------------------------- |
| display | boolean| Yes      | Enable or disable the display of copyright text.                           |
| text    | string | Yes      | Copyright text. Keep it concise. Supported HTML tags: `b`, `i`, `em`, `strong`, `a` |

</details>

###### Example

```json
{
  "room_id": "room01",
  "metadata": {
    "room_title": "Test room",
    "welcome_message": "Welcome to room",
    "room_features": {
      "allow_webcams": true,
      "mute_on_start": false,
      "allow_screen_share": true,
      "admin_only_webcams": false,
      "allow_view_other_webcams": true,
      "allow_view_other_users_list": true,
      "enable_analytics": true,
      "allow_virtual_bg": true,
      "allow_raise_hand": true,
      "auto_gen_user_id": false,
      "recording_features": {
        "is_allow": true,
        "is_allow_cloud": true,
        "is_allow_local": true,
        "enable_auto_cloud_recording": false,
        "only_record_admin_webcams": false,
        "recorder_bot_options": {
          "enable_auto_close_chat_panel": true,
          "duration_after_last_message": 10
        }
      },
      "chat_features": {
        "is_allow": true,
        "is_allow_file_upload": true
      },
      "shared_note_pad_features": {
        "is_allow": true
      },
      "whiteboard_features": {
        "is_allow": true
      },
      "external_media_player_features": {
        "is_allow": true
      },
      "external_broadcasting_features": {
        "is_allow": true,
        "is_allow_rtmp": true,
        "recorder_bot_options": {
          "enable_auto_close_chat_panel": true,
          "duration_after_last_message": 10
        }
      },
      "waiting_room_features": {
        "is_active": false
      },
      "breakout_room_features": {
        "is_allow": true,
        "allowed_number_rooms": 2,
        "allow_return_to_main_room": true,
        "allow_self_select": false,
        "preassigned_rooms": [
          {
            "title": "Team A",
            "user_ids": ["user1", "user2"]
          }
        ]
      },
      "display_external_link_features": {
        "is_allow": true
      },
      "ingress_features": {
        "is_allow": true
      },
      "polls_features": {
        "is_allow": true
      },
      "insights_features": {
        "is_allow": true,
        "transcription_features": {
          "is_allow": true,
          "is_allow_translation": true,
          "is_allow_speech_synthesis": true
        },
        "chat_translation_features": {
          "is_allow": true
        },
        "ai_features": {
          "is_allow": true,
          "ai_text_chat_features": {
            "is_allow": true
          },
          "meeting_summarization_features": {
            "is_allow": true
          }
        }
      },
      "sip_dial_in_features": {
        "is_allow": true,
        "enable_dial_in_on_create": true,
        "hide_phone_number": false
      },
      "end_to_end_encryption_features": {
        "is_enabled": false,
        "enabled_self_insert_encryption_key": false
      }
    },
    "default_lock_settings": {
      "lock_microphone": false,
      "lock_webcam": false,
      "lock_screen_sharing": true,
      "lock_whiteboard": true,
      "lock_shared_notepad": true,
      "lock_chat": false,
      "lock_chat_send_message": false,
      "lock_chat_file_share": false,
      "lock_private_chat": false
    },
    "copyright_conf": {
      "display": true,
      "text": "Powered by <a href=\"https://www.plugnmeet.org\" target=\"_blank\">plugNmeet</a>"
    }
  }
}
```

##### Response

| Field                                           | Type                           | Description               |
|:------------------------------------------------| ------------------------------ | ------------------------- |
| status                                          | boolean                        | Indicates if the request was successful. |
| msg                                             | string                         | Response message.         |
| status_code | number      | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |
| [room_info](https://www.plugnmeet.org/docs/api/room/room-info#room-info) | object | Details about the room.   |

<a id="page-19"></a>

## پایان اتاق

پایان اتاق با room_id درخواست می‌شود. برنامه باید فقط به میزبان مجاز اجازهٔ این عمل بدهد. پایان API و وب‌هوک پایان نشست ممکن است در زمان‌های مختلف دیده شوند و وضعیت نهایی باید با رویداد session_ended و بازیابی API تطبیق داده شود.

این بخش صفحهٔ `api/room/end` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/room/end) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/room/end.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/room/end.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### End room

Endpoint: `/room/endRoom`

This endpoint allows you to forcefully terminate a live room session. When you call this API with a `room_id`, the session is immediately closed, and all participants are disconnected.

This is a server-side administrative action, typically used to moderate sessions or to ensure that a meeting has definitively concluded. For example, you might integrate this into a "End Meeting for All" button in your application's administrative dashboard. Once a room is ended, its data is cleared, and it can be recreated for a new session.

| Field   | Type   | Position | Required | Description |
| ------- | ------ | -------- | :------- | ----------- |
| room_id | string | root     | Yes      | Room Id     |

**Example**:

```json
{
  "room_id": "room01"
}
```

##### Response

| Field  | Type    | Position | Description               |
| :----- | ------- |----------| :------------------------ |
| status | boolean | root     | The status of the request |
| msg    | string  | root     | Response message          |
| status_code | string | root     | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |

<a id="page-20"></a>

## فهرست نشست‌های گذشته

اطلاعات نشست‌های گذشته با صفحه‌بندی دریافت می‌شود. room_id منطقی و sid یک اجرای اتاق را از هم تفکیک کنید. این متد برای تطبیق نشست‌های تمام‌شده و بازسازی وضعیت پس از قطع وب‌هوک مفید است و جای دریافت فایل ضبط را نمی‌گیرد.

این بخش صفحهٔ `api/room/fetch-past` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/room/fetch-past) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/room/fetch-past.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/room/fetch-past.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Fetch Past Rooms Information

Endpoint: `/room/fetchPastRooms`

This endpoint allows you to retrieve historical information about room sessions that have already concluded. After a session ends—either when the last participant leaves or when terminated via the `end` API—its metadata is stored and can be accessed using this API.

This is particularly useful for building features that require a history of past meetings, such as displaying a list of previous sessions in your application's user interface. You can query for one or more `room_id`s to get details like when the session was created, when it ended, and how many participants joined.

##### Request Parameters

| Field    | Type   | Required | Description                                 |
| -------- | ------ | -------- | ------------------------------------------- |
| room_ids | array  | Yes      | Array of room IDs to query.                 |
| from     | number | No       | Starting index for records. Default is 0.   |
| limit    | number | No       | Maximum number of records to return. Default is 20. |
| order_by | string | No       | Sort order: `DESC` or `ASC`. Default is `DESC`. |

**Example Request:**

```json
{
  "room_ids": ["room01"],
  "from": 0,
  "limit": 20,
  "order_by": "DESC"
}
```

##### Response

| Field             | Type                        | Position | Description                       |
| ----------------- | -------------------------- |----------| --------------------------------- |
| status            | boolean                     | root     | Indicates if the request was successful. |
| msg               | string                      | root     | Response message.                 |
| status_code | string | root     | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |
| [result](https://www.plugnmeet.org/docs/api/room/fetch-past#result) | object                      | root     | Contains the results data.        |

###### Result

| Field            | Type                                      | Description                                 |
| ---------------- | ----------------------------------------- | ------------------------------------------- |
| total_rooms      | number                                    | Total number of rooms found for the query.  |
| from             | number                                    | Starting index for the returned records.    |
| limit            | number                                    | Number of records returned.                 |
| order_by         | string                                    | Sort order used for the records.            |
| rooms_list       | Array\<[Past Room Info](https://www.plugnmeet.org/docs/api/room/fetch-past#past-room-info)> | List of past room records.                  |

###### Past Room Info

| Field               | Type   | Description                                                                                   |
| ------------------- | ------ | --------------------------------------------------------------------------------------------- |
| room_title          | string | Title of the room.                                                                            |
| room_id             | string | Unique identifier for the room.                                                               |
| room_sid            | string | SID of the room.                                                                              |
| joined_participants | number | Number of participants who joined (may not be exact; use analytics for detailed information). |
| webhook_url         | string | Webhook URL associated with the room.                                                         |
| created             | string | Room creation time.                                                                           |
| ended               | string | Room end time.                                                                                |
| analytics_file_id   | string | Analytics file identifier, if available.                                                      |

<a id="page-21"></a>

## بررسی فعال‌بودن اتاق

با room_id وضعیت فعال‌بودن اتاق پرسیده می‌شود. نتیجه برای تصمیم ایجاد یا ورود استفاده می‌شود. فعال‌بودن یک مشاهدهٔ لحظه‌ای است و فراخوانی بعدی همچنان باید خطای پایان یا نبود اتاق را مدیریت کند.

این بخش صفحهٔ `api/room/is-active` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/room/is-active) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/room/is-active.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/room/is-active.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Room Status (Active/Inactive)

Endpoint: `/room/isRoomActive`

This endpoint allows you to check whether a room session is currently active. A room is considered "active" if it has been created and has not yet been terminated.

This is a simple way to determine if a meeting is in progress. It's often used in application logic to decide whether to create a new room or allow a user to join an existing one. For example, before calling the `create` API, you can use this endpoint to see if a session with the same `room_id` is already running.

| Field   | Type   | Position | Required | Description           |
| ------- | ------ | -------- | :------- | --------------------- |
| room_id | string | root     | Yes      | The unique ID of the room to check. |

**Example Request:**

```json
{
  "room_id": "room01"
}
```

##### Response

| Field     | Type    | Position | Description                          |
| :-------- | ------- |----------| :----------------------------------- |
| status    | boolean | root     | Indicates if the request was successful. |
| is_active | boolean | root     | Indicates whether the room is currently active. |
| msg       | string  | root     | Response message.                    |
| status_code | string  | root     | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |

<a id="page-22"></a>

## توکن ورود و هویت کاربر

room_id و user_info برای ورود به اتاق فعال ارسال می‌شوند. user_id باید برای هر شرکت‌کننده یکتا و پایدار باشد و is_admin از مجوز واقعی برنامه تعیین شود. با auto_gen_user_id=false ورود مجدد همان شناسه می‌تواند اتصال قبلی را قطع کند. توکن کوتاه‌عمر و یک‌بارمصرف را فوراً در access_token لینک ورود مصرف کنید. نگاشت نام به‌تنهایی برای تشخیص گوینده کافی نیست.

این بخش صفحهٔ `api/room/join` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/room/join) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/room/join.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/room/join.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Join Room

Endpoint: `/room/getJoinToken`

This endpoint generates a temporary token that grants a user access to a specific room. Before diving into the parameters, here are a few key concepts to keep in mind:

*   **Token Lifecycle**: The generated token is short-lived (its validity period is set in your server configuration) and designed for **one-time use**. You should consume it immediately by redirecting the user to the meeting URL. **Do not store this token for future use.**

*   **User ID Uniqueness**:
    *   Plug-N-Meet requires every participant in a session to have a unique `user_id`.
    *   If you created the room with `auto_gen_user_id: false` (the default), you are responsible for providing a unique `user_id` for each user.
    *   If a new user joins with a `user_id` that is already active in the room, the existing participant with that ID will be automatically disconnected. This is useful for allowing users to switch devices seamlessly.

*   **Room Existence**: You can only generate a join token for a room that has already been created and is currently active.

##### Request Parameters

| Field                   | Type   | Required | Description                                 |
| ----------------------- | ------ | -------- | ------------------------------------------- |
| room_id                 | string | Yes      | The ID of the room you want to join.        |
| [user_info](https://www.plugnmeet.org/docs/api/room/join#user-info) | object | Yes      | Information about the user joining the room.|

###### User Info

| Field                           | Type    | Required | Description                                                  |
| --------------------------------| ------- | -------- | ------------------------------------------------------------ |
| name                            | string  | Yes      | The display name of the user.                                |
| user_id                         | string  | Yes      | A unique identifier for the user. **Note:** If the room was created with `auto_gen_user_id: true`, this value will be stored as `ex_user_id` and a random ID will be assigned to the user for the session. |
| is_admin                        | boolean | Yes      | If `true`, the user will join as a moderator with elevated privileges. If `false`, they will join as a standard participant. |
| is_hidden                       | boolean | No       | If `true`, the user will join as a spectator. They will not appear in the participant list and cannot interact. |
| client_type                     | string  | No       | Specifies the client type. Use `"HYBRID_WEB"` for the [hybrid mobile integration](https://www.plugnmeet.org/docs/developer-guide/mobile-app-integration). Defaults to `"WEB"`. |
| [user_metadata](https://www.plugnmeet.org/docs/api/room/join#user-metadata) | object  | Yes      | Additional metadata about the user.                          |

###### User Metadata

| Field                                                           | Type    | Required | Description                                                                                                    |
| --------------------------------------------------------------- | ------- | -------- | -------------------------------------------------------------------------------------------------------------- |
| profile_pic                                                     | string  | No       | URL of the user's profile picture.                                                                             |
| preferred_lang                                                  | string  | No       | Preferred language for the Plug-N-Meet client. See supported values [here](https://github.com/mynaparrot/plugNmeet-client/blob/main/src/helpers/languages.ts). Example: es-ES, bn-BD, de-DE, etc. |
| record_webcam                                                   | boolean | No       | Controls whether this user's webcam is included in server-side recordings. Defaults to `true`. Set to `false` to exclude. |
| ex_user_id                                                      | string  | No       | If empty, the value of `user_id` will be used.                                                                 |
| extra_data                                                      |  map  | No       | A map of custom key-value pairs. Both keys and values must be strings. e.g. `{"key": "value"}`|.                                                                         |
| [lock_settings](https://www.plugnmeet.org/docs/api/room/create#default-lock-settings)    | object  | No       | Lock settings for the user.                                                                                    |

###### Example

```json
{
  "room_id": "room01",
  "user_info": {
    "name": "Your Name",
    "user_id": "Your-Unique-User-Id",
    "is_admin": true,
    "is_hidden": false,
    "client_type": "WEB",
    "user_metadata": {
      "profile_pic": "https://profile.pic/im.jpg",
      "lock_settings": {
        "lock_microphone": false,
        "lock_webcam": false,
        "lock_screen_sharing": true,
        "lock_chat": false,
        "lock_chat_send_message": false,
        "lock_chat_file_share": false
      }
    }
  }
}
```

##### Response

| Field  | Type    | Position | Description               |
| ------ |---------|----------| ------------------------- |
| status | boolean | root     | Indicates if the request was successful. |
| msg    | string  | root     | Response message.         |
| status_code | string | root     | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |
| token  | string  | root     | The join token.           |

###### Using the Join Token

Once you receive the `token`, you can grant the user access to the meeting room. There are two ways to pass the token to the Plug-N-Meet client:

1.  **As a URL parameter (recommended)**: Pass the token in the query string as `access_token`.
2.  **As a cookie**: Set a cookie named `pnm_access_token` with the token as its value.

The URL parameter is the most straightforward method. Here's an example:

```
https://Your-Plug-N-Meet-Server.com/?access_token=<TOKEN HERE>
```

If you are hosting the [plugNmeet-client](https://github.com/mynaparrot/plugNmeet-client) static files elsewhere, you should also pass the token in the same way. For example, if your static files are in a `conference` directory:

```
https://Your-Domain.com/conference/?access_token=<TOKEN HERE>
```

**Note:** You can also use the [getClientFiles](https://www.plugnmeet.org/docs/api/get-client-files) API to retrieve all the CSS and JS files required to display the interface manually. This allows you to embed the Plug-N-Meet client anywhere without worrying about static build files.

##### Custom Design

You can add an extra query parameter called `custom_design` to the join link to apply a unique design for each user. The value of `custom_design` should be a URL-encoded JSON string.

For a full list of supported parameters and their descriptions, please refer to the **[Design Parameters](https://www.plugnmeet.org/docs/developer-guide/design-customisation#design-parameters)** section in our Developer Guide.

**Example:**

```
https://Your-Domain.com/conference/?access_token=<TOKEN HERE>&custom_design=%7B%22primary_color%22%3A%22%23004D90%22%2C%22secondary_color%22%3A%22%24AEF7%22%7D
```

<a id="page-23"></a>

## جزئیات اتاق فعال

این متد اطلاعات یک اتاق فعال و شرکت‌کنندگان آن را برمی‌گرداند. اطلاعات نشست و کاربر برای داشبورد و تطبیق وضعیت مفید است. نتیجهٔ زنده باید در کنار شناسهٔ اجرای جلسه ثبت شود تا اتاق‌های تکرارشده با هم مخلوط نشوند.

این بخش صفحهٔ `api/room/room-info` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/room/room-info) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/room/room-info.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/room/room-info.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Get Active Room Information

Endpoint: `/room/getActiveRoomInfo`

This endpoint retrieves detailed, real-time information about a specific active room. While `isRoomActive` simply tells you if a room is running, this API provides a comprehensive snapshot of the session, including its current state, a list of all participants, and their individual details.

This is useful for building administrative dashboards or moderation tools where you need to monitor live meetings. You can see who is in the room, when they joined, and other metadata associated with both the room and the participants.

##### Request Parameters

| Field   | Type   | Position | Required | Description           |
| ------- | ------ | -------- | :------- | --------------------- |
| room_id | string | root     | Yes      | The unique ID of the room. |

**Example Request:**

```json
{
  "room_id": "room01"
}
```

##### Response

| Field   | Type                  | Position | Description               |
| ------- | --------------------- | -------- | ------------------------- |
| status  | boolean               | root     | Indicates if the request was successful. |
| msg     | string                | root     | Response message.         |
| status_code | string | root     | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |
| room    | object\<[Room](https://www.plugnmeet.org/docs/api/room/room-info#room)>| root     | Contains room details.    |

###### Room

| Field                | Type                                         | Description                       |
| -------------------- | -------------------------------------------- | --------------------------------- |
| room_info            | object\<[Room Info](https://www.plugnmeet.org/docs/api/room/room-info#room-info)>             | Details about the room.           |
| participants_info    | array\<[Participant Info](https://www.plugnmeet.org/docs/api/room/room-info#participant-info)>| List of current participants.     |

###### Room Info

| Field               | Type    | Description                                 |
| ------------------- | ------- | ------------------------------------------- |
| room_title          | string  | Title of the meeting.                       |
| room_id             | string  | Unique room ID.                             |
| sid                 | string  | Room unique session ID.                                   |
| joined_participants | number  | Total number of users who have joined.      |
| is_running          | boolean | Indicates if the room is currently active.  |
| is_recording        | boolean | Indicates if recording is in progress.      |
| is_active_rtmp      | boolean | Indicates if RTMP streaming is active.      |
| creation_time       | number  | Room creation time (Unix timestamp).        |
| metadata            | string  | Room metadata.                              |
| webhook_url         | string  | Webhook URL associated with the room.       |

###### Participant Info

| Field      | Type    | Description                                 |
| ---------- | ------- | ------------------------------------------- |
| sid        | string  | Participant unique session ID.                            |
| identity   | string  | Participant user ID.                        |
| name       | string  | Participant name.                           |
| state      | string  | Participant state.                          |
| metadata   | string  | Participant metadata.                       |
| joined_at  | number  | Time the participant joined (Unix timestamp).|
| version    | number  | Version information.                        |

<a id="page-24"></a>

## فهرست اتاق‌های فعال

این متد نمای کلی اتاق‌های فعال سرور را می‌دهد. برای پایش و reconciliation کاربرد دارد. فراخوانی مدیریتی از backend انجام شود و نمایش اتاق‌ها در آی‌نشست با مجوز سازمان و کاربر محدود شود.

این بخش صفحهٔ `api/room/rooms-info` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/room/rooms-info) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/room/rooms-info.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/room/rooms-info.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Get active rooms info

Endpoint: `/room/getActiveRoomsInfo`

This endpoint provides a comprehensive overview of all currently active rooms on the server. It returns an array of room objects, each containing the same detailed information as the `/room/getActiveRoomInfo` endpoint, including participant lists and room metadata.

This is a powerful tool for server-wide monitoring and administration. It's ideal for building a real-time dashboard that shows all ongoing meetings, the number of participants in each, and other session details.

**Example**:

```json
{}
```

##### Response

| Field  | Type    | Position | Description                                                                               |
| :----- | ------- | -------- | :---------------------------------------------------------------------------------------- |
| status | boolean | root     | The status of the request                                                                 |
| msg    | string  | root     | Response message                                                                          |
| status_code | string | root     | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |
| rooms  | array\<[room](https://www.plugnmeet.org/docs/api/room/room-info#room)>   | root     | Array of Room |

<a id="page-25"></a>

## ارسال فایل به وایت‌برد زنده

اتاق باید فعال باشد و دست‌کم یک presenter داخل آن حضور داشته باشد. درخواست multipart فقط یکی از document یا document_link دارد و Room-Id الزامی است. امضا از مقدار Room-Id ساخته می‌شود. محدودیت حجم از max_size_whiteboard_file می‌آید. طبق این نسخهٔ مستندات، فایل‌های اداری و PDF پشتیبانی می‌شوند و تصویر پذیرفته نمی‌شود. آپلود هم‌زمان همان اتاق خطای 409 می‌دهد. تبدیل پس از timeout درخواست ممکن است در پس‌زمینه ادامه یابد.

این بخش صفحهٔ `api/room/upload-whiteboard-file` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/api/room/upload-whiteboard-file) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/api/room/upload-whiteboard-file.md) · [نسخهٔ محلی کامل](docs/upstream/docs/api/room/upload-whiteboard-file.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Upload Whiteboard File

Endpoint: `/room/uploadWhiteboardFile`

While presenters can upload files directly through the Plug-N-Meet client interface, this API endpoint provides the flexibility to inject files into a live session from an external system. It acts as a bridge between your application's backend and the Plug-N-Meet room, allowing for powerful, automated workflows.

For example, you could build an integration that allows users to select a file from a cloud storage service like Google Drive or an S3 bucket. Your application can either download the file and then push it to the server using a `multipart/form-data` request, or you can provide a direct download URL using the `document_link` parameter. If you use `document_link`, it is crucial that you provide the final, direct download link, as the server will not perform any authentication. Your application is responsible for handling any necessary authentication with the storage service and generating a publicly accessible link.

After a successful upload via this API:
1.  The presenter receives a notification within the room.
2.  The new file appears in the **whiteboard files list**.
3.  The presenter can then choose when to display the file on the whiteboard for all participants.

This allows for pre-loading multiple files into the session's file list. Note that each API request can only contain one file at a time.

##### Prerequisites

For this API call to succeed, the following conditions must be met:
1.  The session (`room_id`) must be active.
2.  At least one user with a `presenter` role must be present in the room.

This functionality is different from the `preload_file` option available in the [Create Room API](https://www.plugnmeet.org/docs/api/room/create). The `preload_file` option accepts a direct URL and prepares the file when the room is created (before it is active). This endpoint, in contrast, is designed for uploading files to an already running session.

##### Request

The request must be sent as a `multipart/form-data` request and include the required headers.

###### Headers

| Header           | Type   | Required | Description                                                                 |
| ---------------- | ------ | -------- | --------------------------------------------------------------------------- |
| `API-KEY`        | string | Yes      | Your API key for authentication.                                            |
| `HASH-SIGNATURE` | string | Yes      | The HMAC-SHA256 signature. See [Authentication](https://www.plugnmeet.org/docs/api/room/upload-whiteboard-file#authentication-for-multipart-requests) for details. |
| `Room-Id`        | string | Yes      | The unique identifier of the active room to which you want to upload the file. |

###### Body

You must provide either `document` (for direct file upload) or `document_link` (for server-side download), but not both. All uploads are subject to the [File Constraints](https://www.plugnmeet.org/docs/api/room/upload-whiteboard-file#file-constraints) detailed below.

| Field          | Type   | Required | Description                                                                 |
| -------------- | ------ | -------- |-----------------------------------------------------------------------------|
| `document`     | file   | No       | The file to be uploaded. See [File Constraints](https://www.plugnmeet.org/docs/api/room/upload-whiteboard-file#file-constraints).         |
| `document_link`| string | No       | A publicly accessible URL from which the server will download the file. **Important**: The server will not perform any authentication, so your application must provide a final, direct download link. See File Constraints. |

###### File Constraints

-   **Supported File Types**: This endpoint is designed for office documents and presentations. Supported formats include **PDF, DOC, DOCX, PPT, PPTX, etc.** Image files (e.g., JPG, PNG) are not supported and will be rejected.
-   **File Size Limit**: The size of the file, whether uploaded directly or downloaded via `document_link`, must not exceed the `max_size_whiteboard_file` limit defined in your server's configuration.

##### Authentication for Multipart Requests

Unlike standard JSON requests, `multipart/form-data` requests require a special authentication method. For this specific endpoint, the `HASH-SIGNATURE` **must be generated from the value of the `Room-Id` header**.

-   **`Content-Type`**: Should be `multipart/form-data`. cURL and most HTTP clients will set this automatically when you use form fields.
-   **`HASH-SIGNATURE`**: An HMAC-SHA256 signature generated using your API Secret, with the `Room-Id` as the message body.

##### Example cURL Requests

This section demonstrates how to correctly generate the signature and send a file upload request using both methods.

###### Example 1: Uploading a file directly

```bash
# Your API credentials
API_KEY="plugnmeet"
SECRET="EXAMPLE_API_SECRET"
ROOM_ID="room01"

# 1. Generate the signature from the Room-Id value.
SIGNATURE=$(echo -n "$ROOM_ID" | openssl dgst -sha256 -mac HMAC -macopt key:"$SECRET" | awk '{print $2}')

# 2. Make the POST request with the correct headers and form data.
curl -X POST 'https://plugnmeet.example.com/auth/room/uploadWhiteboardFile' \
--header "API-KEY: $API_KEY" \
--header "HASH-SIGNATURE: $SIGNATURE" \
--header "Room-Id: $ROOM_ID" \
--form 'document=@"/path/to/your/presentation.pdf"'
```

###### Example 2: Uploading a file via a link

```bash
# Your API credentials
API_KEY="plugnmeet"
SECRET="EXAMPLE_API_SECRET"
ROOM_ID="room01"
DOCUMENT_LINK="https://example.com/my_presentation.pdf"

# 1. Generate the signature from the Room-Id value.
SIGNATURE=$(echo -n "$ROOM_ID" | openssl dgst -sha256 -mac HMAC -macopt key:"$SECRET" | awk '{print $2}')

# 2. Make the POST request with the correct headers and form data.
curl -X POST 'https://plugnmeet.example.com/auth/room/uploadWhiteboardFile' \
--header "API-KEY: $API_KEY" \
--header "HASH-SIGNATURE: $SIGNATURE" \
--header "Room-Id: $ROOM_ID" \
--form "document_link=$DOCUMENT_LINK"
```

##### Response

**Important Notes on Concurrent Uploads and Timeouts:**

The server enforces a strict one-file-at-a-time upload policy per room to ensure data integrity. This applies whether you are uploading a file directly (`document`) or providing a URL (`document_link`).

-   **Concurrent Uploads**: The server will reject any new upload requests for a room if another upload is already in progress. If you attempt a concurrent upload, the API will respond with an HTTP status code of `409` (Conflict).
-   **Timeouts and Background Processing**: The file conversion process can be resource-intensive. If the initial request times out, the server will continue processing the file in the background. Once completed, it will notify the client, and the room's upload lock will be released, allowing for another file to be uploaded.
-   **Upload Lock Duration**: By default, the upload lock for a room is set to 5 minutes. If the file processing does not complete within this window, the operation will time out. It is crucial to select files that can be reasonably processed within this timeframe.

If the file is uploaded and processed successfully, the server will broadcast the necessary information to all participants in the room, making the new file available on the whiteboard.

| Field       | Type    | Description                              |
| ----------- | ------- | ---------------------------------------- |
| status      | boolean | Indicates if the request was successful. |
| msg         | string  | Response message.                        |
| status_code | string  | Response [status code](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common_api.proto#L10). |

<a id="page-26"></a>

## ظرفیت و روش آزمون

این صفحه نتایج لپ‌تاپ، یک سرور و خوشه را با سناریوهای متفاوت فعالیت و رسانه گزارش می‌کند. تعداد کاربر به‌تنهایی ظرفیت را مشخص نمی‌کند و ویدئو، پهنای باند و تعداد ضبط بر نتیجه اثر می‌گذارند. اعداد صفحه برای همان محیط و روش آزمون معنا دارند. در آی‌نشست مبنای پذیرش اولیه ۲ اتاق، مجموعاً ۱۰ کاربر و ۲ ضبط هم‌زمان است و باید سلامت رسانه و فایل خروجی همراه با CPU، RAM، شبکه و دیسک سنجیده شود.

این بخش صفحهٔ `benchmark` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/benchmark) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/benchmark.md) · [نسخهٔ محلی کامل](docs/upstream/docs/benchmark.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### How many people can plugNmeet handle? We tested it.

> Tested September 2026 · plugNmeet [server](https://github.com/mynaparrot/plugNmeet-server) v2.5.1 · [client](https://github.com/mynaparrot/plugNmeet-client) v2.5.2 · [loadtester](https://github.com/mynaparrot/plugNmeet-loadtester) v1.0.1

We often say plugNmeet is lightweight and scalable. That is easy to claim, harder to prove — so we measured it.

We filled real rooms with simulated participants that behave like browser users: they join, chat, draw on the whiteboard, raise hands, react, and send heartbeats. Then we watched what happened on three setups you actually care about:

1. **One laptop doing everything** — the worst case.
2. **One normal server installed with our [single install script](https://www.plugnmeet.org/docs/installation)** — what most self-hosters run.
3. **A production cluster** — what you grow into.

If you are not technical, the next section gives you the whole story in 60 seconds. If you are, keep scrolling for tables, latencies, and how to reproduce it.

##### If you only read 60 seconds

Think of a meeting as two jobs: **talking about the meeting** (who joined, chat, whiteboard, reactions) and **moving video and audio**. The first job is cheap for plugNmeet. The second depends on how many cameras are on at once.

Here is what we found:

- **A small class or team meeting (20–100 people in one room): easy.** Even a laptop handles it. A single server answers in ~25–30 ms with zero drops.
- **A large webinar (300–500 people in one room, 1–2 cameras on, everyone else watching): comfortable on one server.** We ran 300 and 500 users in a single room on a 4-vCore / 16 GB server with **zero missed heartbeats, zero errors, ~100% delivery**. Adding live video did not break anything: video joined in ~350 ms and traveled in ~25 ms.
- **Many small rooms are quieter than one giant room.** 10 rooms × 30 users and 5 rooms × 100 users both ran perfectly on the same single server, with slightly better response times than putting all 500 in one room.
- **Around 900 people in one room is where a single 4-vCore server stops.** At 1,000 attempted joins, 901 got in, 99 timed out, and heartbeats started missing. The server did not crash — it just ran out of headroom.
- **1,000 people in one interactive room is possible, but not on one box.** On our managed cluster it ran clean for 18 minutes with zero drops and no lag for two real humans who joined in normal browsers. That cluster uses a custom multi-server media routing setup that is **not part of the open-source single-server install** (more on that below).

Bottom line: **start with one server, grow only when you must.** Most schools, creators, and teams never need more than that. When you do, follow the [Scalable Deployment guide](https://www.plugnmeet.org/docs/developer-guide/scalable-setup).

##### How to read this page

- **Non-technical track:** read the summaries, the “what this means for you” section, and skip the tables.
- **Technical track:** read everything, especially [Methodology](https://www.plugnmeet.org/docs/benchmark#methodology-how-we-tested) and [Test it yourself](https://www.plugnmeet.org/docs/benchmark#test-it-yourself).

Three terms appear everywhere:

- **p50 / p95:** sort every response time from fastest to slowest. **p50** is the middle — half were faster. **p95** is what 95 out of 100 beat — it shows the unlucky few. Under ~100 ms feels instant.
- **Missed heartbeats:** every user asks “still there?” every 10 seconds. Missed = 0 means nobody lost connection. Anything above 0 means trouble.
- **Fan-out:** when one person chats or draws, the server must deliver it to everyone else. 100% means nobody missed anything.

A run counts as clean when: `joined == requested`, errors show `none`, missed heartbeats = 0, reconnects = 0, fan-out ≈ 100%.

##### The three setups

###### 1. Laptop (worst case)

One Lenovo laptop (i9-13900H, 32 GB RAM) ran the plugNmeet server, NATS message bus, database, **and** all simulated users at once. Deliberately unfair — if it works here, it works on real hardware.

###### 2. Single server via installer (the one most readers will run)

This is the missing piece people asked for: a plain server set up exactly the way the docs describe, using the [Installation Quick Start](https://www.plugnmeet.org/docs/installation).

- **plugNmeet server:** OVH `b3-16-flex`, 4 vCore, 16 GB RAM, Gravelines, France (GRA11). Full install: plugNmeet server + LiveKit + NATS + Redis + MariaDB + HAProxy on one box. Recorder installed but **not recording** during tests.
- **Load generator:** separate OVH box, 8 vCore, 32 GB RAM, Warsaw, Poland (WAW1), running our open-source [plugnmeet-loadtester](https://github.com/mynaparrot/plugnmeet-loadtester).
- **Why two cities?** To include real internet latency. Baseline round-trip between tester and server was ~25–30 ms — similar to what a European audience would see. This is not a localhost trick.
- **Join speed:** 1 user per second globally (so 300 users need at least 5 minutes; we used 10-minute runs, 2 minutes for tiny rooms, 18 minutes for the 1,000-user push).

No special tuning. Firewall as per [Firewall guide](https://www.plugnmeet.org/docs/firewall). Fresh room ID per run.

###### 3. Production cluster

Our own cloud edge cluster across two zones, public entry points only. Same application code as open source, but media routing uses a **custom multi-LiveKit setup** that spreads one room across multiple media servers. The open-source install runs one room on a single LiveKit node — that distinction matters for the 1,000-user headline below. Resource usage stayed low throughout — the point of the cluster is headroom and distribution, not struggling hardware.

##### Results: laptop

| Users | Joined | Errors | Missed | Heartbeat p50 / p95 | Laptop CPU | Server process |
|-------|--------|--------|--------|---------------------|------------|----------------|
| 20 | 20/20 | none | 0 | 6.0 / 9.6 ms | 8% | under 1% of one core |
| 100 | 100/100 | none | 0 | 5.1 / 9.7 ms | 14% | 1.4% |
| 300 | 300/300 | none | 0 | 5.2 / 18.3 ms | 37% | 3.7% |
| 500 | 500/500 | none | 0 | 4.4 / 27.9 ms | 53% | 6.0% |
| 1000 | 795/1000 | 146 | 673 | 9.2 / 2,719 ms | 59% | 7.9% |

Translation: **clean through 500**, and at 1,000 the laptop itself gave out while the server process stayed under 8% of one core. The bottleneck was the test bench, not plugNmeet.

##### Results: single server, meeting activity only (no video yet)

Same webinar shape — chat, whiteboard, notepad, reactions, hands, heartbeats — but no cameras. This isolates the “talking about the meeting” cost.

| Users in one room | Joined | Missed | Heartbeat p50 / p95 | Chat delivery p50 / p95 | Whiteboard pointer p50 / p95 | Fan-out |
|-------------------|--------|--------|---------------------|-------------------------|------------------------------|---------|
| 20 | 20/20 | 0 | 27.8 / 53.1 ms | 26.3 / 26.5 ms | 26.1 / 26.5 ms | 100% |
| 50 | 50/50 | 0 | 27.4 / 53.5 ms | 26.2 / 26.9 ms | 26.1 / 26.8 ms | 100% |
| 100 | 100/100 | 0 | 28.0 / 52.8 ms | 26.0 / 27.4 ms | 26.1 / 27.6 ms | 100% |
| 300 | 300/300 | 0 | 27.9 / 44.2 ms | 26.5 / 33.3 ms | 26.2 / 33.2 ms | 100% |
| 500 | 500/500 | 0 | 33.8 / 103.1 ms | 27.7 / 43.4 ms | 27.8 / 44.4 ms | ~100% (chat 99.9%, rest 100%) |
| 1000 | 901/1000 | 8,109 | 163.8 / 7,734 ms | 73.2 / 224.5 ms | 71.3 / 188.4 ms | degraded |

Notes:

- 20–100 ran for 2 minutes each, 300 and 500 for 10 minutes each, 1,000 for 18 minutes.
- At 500, everything still joined and heartbeats never missed, but p95s start climbing (heartbeat p95 ~103 ms, user-list fetch p50 ~73 ms, p95 ~180 ms). That is the server telling you it is getting busy, not failing.
- At 1,000, 99 users failed with `initialData.timeout` and the rest saw seconds-long replies. Honest ceiling on this 4-vCore box: **somewhere around 900 in one room for meeting activity alone**.
- Server-side during these runs: NATS, LiveKit, and HAProxy were the busiest processes (bursts around 1–1.6 cores for the proxy/bus during fan-out), none saturated the 4 cores. The plugNmeet app process itself stayed light — the heavy lifting is copying each message to everyone, which is NATS’s job.

##### Results: same single server, now with live video

We then turned cameras on using `--media video`. Notation: `publishers=2v+5a` means 2 users publish camera+mic and 5 publish mic-only; `subscribers=100/300` means 100 of the 300 users connect to LiveKit, the other 200 are meeting-activity-only.

| Room | Video setup | Joined | Missed | Core (chat/pointer p50) | Video join p50 / p95 | Video travel p50 / p95 | Fan-out |
|------|-------------|--------|--------|-------------------------|----------------------|------------------------|---------|
| 50 users | 1 camera, all 50 watching | 50/50 | 0 | ~26 ms | 349 / 450 ms | 25.2 / 25.7 ms | 100% |
| 20 users | 10 cameras, all 20 watching | 20/20 | 0 | ~26 ms | 350 / 450 ms | 25.2 / 25.4 ms | 100% |
| 300 users | 2 cameras, all 300 watching | 300/300 | 0 | ~26 ms | 357 / 395 ms | 25.2 / 27.8 ms | 100% |
| 300 users (mixed) | 2 cameras + 5 mics, 100 on video, 200 activity-only | 300/300 | 0 | ~26 ms | 351 / 448 ms | 25.2 / 27.6 ms | 100% |

What this means in plain words:

- **Viewers are cheap, publishers cost.** Going from 1 to 10 publishers, or from 50 to 300 viewers, did not move core latencies. Video join stayed ~350 ms, video travel stayed ~25 ms.
- **About the TURN warnings in the video logs.** On the single-server runs you will see `Fail to refresh permissions` / `write: broken pipe` to port 443. That is the installer’s co-located TURN (sharing port 443 TCP on the same box) struggling with refreshes from a distant tester over TCP. It did not affect these results (0 missed, 100% fan-out). On our cluster, which uses a dedicated TURN server, we do not see these errors — careful clustering with TURN separated solves most of it. For production video, keep UDP 50000–60000 open per the [Firewall guide](https://www.plugnmeet.org/docs/firewall) and separate TURN/media as in [Scalable Deployment](https://www.plugnmeet.org/docs/developer-guide/scalable-setup).
- We did **not** test 50+ simultaneous publishers here. Each live camera multiplies media work, so treat the rows above as “webinar shape” (a few presenters, many watchers), not “everyone on camera.”

##### Results: same single server, many small rooms

Same 300 or 500 total users, but split up. Real schools and SaaS products look like this, not one giant room.

| Layout | Joined | Missed | Heartbeat p50 / p95 | Fan-out |
|--------|--------|--------|---------------------|---------|
| 10 rooms × 30 users (300 total) | 300/300 | 0 | 32.7 / 58.0 ms | 100% |
| 5 rooms × 100 users (500 total) | 500/500 | 0 | 29.2 / 77.9 ms | 100% |

Splitting helps: per-message fan-out is smaller, lists are shorter, p95s stay flatter than putting all 500 in one room. If your use case is classes or team meetings, this is the row to quote.

##### Results: cluster, 1,000 in one room

For completeness, the headline from our managed cluster (custom multi-media-server routing, same app code otherwise):

- Single room, **2 live webcams**, everyone watching, whiteboard + chat + reactions + hands active, **18 minutes**.
- **0 missed heartbeats** (58,161/58,161), **0 reconnects**, **100% fan-out** on every message kind.
- Whiteboard p50 ~15 ms, pointer ~29 ms, chat ~29 ms, heartbeat ~68 ms, video round-trip ~5 ms across 114,000+ samples.
- App server ~21% of one core per node. Each client absorbed ~78 messages/sec sustained.
- Two real humans joined in normal browsers: **no perceptible lag**.

Do not compare this directly to the single-server 1,000-user row. The difference is the media layer: the cluster spreads one room over multiple LiveKit backends; the open-source install keeps one room on one LiveKit node. Everything else (NATS design, API, recorder split) follows the [Scalable Deployment guide](https://www.plugnmeet.org/docs/developer-guide/scalable-setup).

##### What this means for you

- **20–100 in a room:** any small VPS from the [Installation guide](https://www.plugnmeet.org/docs/installation) works. Test with 2 vCPU / 4 GB and grow if needed.
- **300–500 in one room (webinar shape):** a single 4-vCore / 16 GB server is comfortable. Keep the recorder on a separate box if you record a lot — transcoding steals CPU from live meetings (see [Scalable Deployment](https://www.plugnmeet.org/docs/developer-guide/scalable-setup)).
- **500+ total across rooms:** same single server still works (5 × 100 was clean). Scale vertically first before clustering.
- **~900+ in one interactive room:** plan a cluster, or use the broadcast pattern: private studio + RTMP to YouTube/Facebook, bring guests on stage via waiting room. See [Hosting large-scale events](https://www.plugnmeet.org/blog/hosting-large-scale-events-the-smart-way).
- **Video rule of thumb:** viewers add little; each extra camera/mic adds real media work. Place media servers close to participants — geography matters more than CPU for video quality.
- **Bandwidth matters more than RAM.** Our install docs recommend minimum 100 Mbits/sec; more is always better for video.

##### Honest limits: what we did not test

- **Recorder ON vs OFF:** all rows above ran with the recorder idle. Recording + transcoding on the same box will lower capacity. Isolate it first.
- **Many publishers, 1,000 with video, slow/mobile networks, multi-hour events:** not covered here. They are good “run it yourself” tests — the tool supports `--video-publishers`, `--audio-publishers`, `--subscribers`, `--rooms`, and `--disable` to shape them.
- **Open-source single-room media ceiling:** open-source LiveKit keeps one room on one media node. That is why the single-server 1,000-user push stopped around 900 while the custom cluster went further. If huge single rooms are your product, budget for media clustering work.

If you need one of these gaps filled with numbers before you commit, tell us — we append new results to this page as they come in.

##### Methodology: how we tested

- **Tool:** [plugnmeet-loadtester](https://github.com/mynaparrot/plugnmeet-loadtester). Simulated users speak the same HTTP + NATS-WebSocket protocol as the web client, so the server cannot tell them apart. Optional `--media` flag drives LiveKit with real black-frame video.
- **Behavior per user:** join handshake (`verifyToken`, `getJoinToken`, `initialData`), heartbeat every 10 s, online/user lists, chat, whiteboard pointer + scene updates, notepad sync, reactions, hands. Whiteboard cursors dominate (~80%+ of messages) — real audiences are quieter, so these results are a lower bound, not a ceiling.
- **Commands (examples):**
  ```bash
  # meeting activity only, 500 in one room
  ./bin/loadtest --server https://your-server.example.com \
    --api-key YOUR_KEY --api-secret YOUR_SECRET \
    --room bench-01 --users 500 --duration 10m

  # webinar shape: 300 users, 2 cameras, everyone watches
  ./bin/loadtest --server https://your-server.example.com \
    --api-key YOUR_KEY --api-secret YOUR_SECRET \
    --room bench-video --users 300 --duration 10m \
    --media video --video-publishers 2 --subscribers 300

  # mixed: 300 users, 2 cameras + 5 mics, 100 on video
  ./bin/loadtest --server https://your-server.example.com \
    --api-key YOUR_KEY --api-secret YOUR_SECRET \
    --room bench-mixed --users 300 --duration 10m \
    --media video --video-publishers 2 --audio-publishers 5 --subscribers 100

  # many rooms: 5 rooms x 100 users
  ./bin/loadtest --server https://your-server.example.com \
    --api-key YOUR_KEY --api-secret YOUR_SECRET \
    --room bench-multi --rooms 5 --users 100 --duration 10m
  ```
- **Rules we followed:** fresh `--room` per run (stale rooms block joins), `--join-rate 1/s` default, size the generator at ~10 MB RAM per simulated user, only test servers you own.
- **What we recorded:** `connect`, join paperwork, `ping`, `delivery.*` RTT + fan-out, `mediaJoin`/`mediaRtt` when video is on, plus generator CPU/RSS and machine load. Full flag reference is in the [tool README](https://github.com/mynaparrot/plugnmeet-loadtester).

##### Test it yourself

You do not need our hardware. Install with the [Installation Quick Start](https://www.plugnmeet.org/docs/installation), open [Firewall ports](https://www.plugnmeet.org/docs/firewall), then run the commands above starting at `--users 50` and doubling. Stop when missed heartbeats rise above 0 or fan-out drops below ~100% — that is your server’s honest limit for that shape.

###### Numbers alone are not enough — join with real browsers

The load tester tells you the server stayed up. Only a human can tell you it *felt* good. So do what we did: mix bots with at least **2 real users in normal browsers** in the same room while the run is going.

- Bots provide the stress (hundreds of joins, chat, whiteboard flood).
- Real users check what bots cannot: whiteboard sync with your own eyes, video smoothness, audio clarity, and whether clicks feel instant.

Just open the room URL twice (two laptops, or one laptop + one phone) and use it normally while the bots run. If the whiteboard stays in sync and audio/video stays smooth, the numbers and the experience agree.

###### Test from a different zone if you can

Run the tester and join the real browsers from a different city or network than the server — that adds real internet latency, like your actual audience. Our single-server runs did exactly this (server in France, tester in Poland, ~25–30 ms baseline).

One honest caveat: our managed cluster adds **GEO routing** (users connect to the nearest media region), which is not part of the open-source single-server install. So expect higher video latency at distance on a single server — that is geography, not a bug. Place media servers close to participants when it matters.

Related reading: [Scalable Deployment](https://www.plugnmeet.org/docs/developer-guide/scalable-setup) · [Installation](https://www.plugnmeet.org/docs/installation) · [Backend Architecture Deep Dive](https://www.plugnmeet.org/blog/backend-architecture-deep-dive) · [Scaling Architecture Saves Money](https://www.plugnmeet.org/blog/scaling-architecture-saves-money) · [Hosting large-scale events](https://www.plugnmeet.org/blog/hosting-large-scale-events-the-smart-way)

<a id="page-27"></a>

## شخصی‌سازی ظاهر

پیکربندی config.js و پارامترهای طراحی امکان تغییر رنگ، لوگو و اجزای رابط را می‌دهند. custom_design در لینک ورود باید JSON معتبر و URL-encoded باشد. برای تغییرات فراتر از پارامترهای رسمی، CSS یا سورس کلاینت استفاده می‌شود. تغییر ظاهر ابتدا با نسخهٔ فعلی کلاینت و روی میزبان و مهمان آزموده شود.

این بخش صفحهٔ `developer-guide/design-customisation` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/developer-guide/design-customisation) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/developer-guide/design-customisation.md) · [نسخهٔ محلی کامل](docs/upstream/docs/developer-guide/design-customisation.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Design Customization

This guide explains how to customize the look and feel of the plugNmeet user interface. You can achieve a fully branded, white-label experience by using a combination of simple configuration parameters and your own custom CSS.

##### Method 1: Using `config.js` for Basic Customization

For quick and easy branding, you can edit the `config.js` file located in the `dist/assets/` directory of the `plugNmeet-client` package.

In the latest version, all configuration is handled through the `window.plugNmeetConfig` object. You can define your design customizations within the `designCustomization` property of this object.

> **Note:** For backward compatibility, the legacy `window.DESIGN_CUSTOMIZATION` variable is still supported, but we recommend using the new `window.plugNmeetConfig` structure for all new deployments.

###### Configuration Example

Here is how to set up your `config.js` with design customizations:

```javascript
window.plugNmeetConfig = {
  // ... other config options ...

  // Design customization in JS Object or JSON format.
  designCustomization: {
    primary_color: "#004D90",
    primary_btn_bg_color: "#00a1f28c",
    primary_btn_text_color: "#ffffff",
    secondary_color: "#24AEF7",
    secondary_btn_bg_color: "#ffffff8c",
    secondary_btn_text_color: "#0c131a",
    header_bg_color: "#45b3ec",
    footer_bg_color: "#45b3ec",
    footer_icon_bg_color: "#004d90",
    footer_icon_color: "#ffffff",
    side_panel_bg_color: "#04a2f3",
    background_color: "#0b7db4",
    background_image: "https://mydomain.com/custom_bg.png",
    custom_css_url: "https://mydomain.com/plugNmeet_desing.css",
    custom_logo: "https://mydomain.com/logo.png"
  }
};
```

**Dynamic Customization**
You can also pass this same JSON object as a URL parameter named `custom_design` when joining a room. This allows you to apply different designs for different users or meetings dynamically.


###### Design Parameters

> **Important:** All external resources (`background_image`, `custom_css_url`, `custom_logo`, etc.) must be served over **HTTPS**. Using non-HTTPS (`http://`) URLs will cause mixed-content errors in the browser and prevent the resources from loading.

**Dark Mode Limitation**
Most **color-related parameters** (e.g., `primary_color`, `header_bg_color`) currently apply only to the **Light Theme**. If a user switches to Dark Mode, these custom colors will be overridden by the default dark theme styles.

However, global assets like `custom_logo`, `background_image`, and `background_color` will still apply. For full control over Dark Mode styling, we recommend using the `custom_css_url` option to provide your own advanced CSS overrides.


| Field                      | Description                                                                                                                                     |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| `primary_color`            | The primary brand color for buttons, icons, and other key highlights. (Hex format)                                                              |
| `primary_btn_bg_color`     | The background color for primary buttons. (Hex format)                                                                                          |
| `primary_btn_text_color`   | The text color for primary buttons. (Hex format)                                                                                                |
| `secondary_color`          | The secondary brand color, often used for accents. (Hex format)                                                                                 |
| `secondary_btn_bg_color`   | The background color for secondary buttons. (Hex format)                                                                                        |
| `secondary_btn_text_color` | The text color for secondary buttons. (Hex format)                                                                                              |
| `header_bg_color`          | The background color for the main header bar. (Hex format)                                                                                      |
| `footer_bg_color`          | The background color for the main footer and control bar. (Hex format)                                                                          |
| `footer_icon_bg_color`     | The background color for the icons in the footer. (Hex format)                                                                                  |
| `footer_icon_color`        | The color of the icons in the footer. (Hex format)                                                                                              |
| `side_panel_bg_color`      | The background color for the side panel (e.g., chat, participants etc). (Hex format)                                                            |
| `background_color`         | Sets the default background color of the main meeting room. (Hex format)                                                                        |
| `background_image`         | Sets a custom background image for the room. This must be a direct HTTPS URL. If `background_image` is set, `background_color` will be ignored. |
| `custom_css_url`           | A URL to your own custom CSS file for advanced, granular customization. This must be a direct HTTPS URL. See Method 2 below.                    |
| `custom_logo`              | A URL for your brand's logo, displayed on the welcome screen and in the header. **Must be a direct HTTPS URL.**                                 |

---

##### Method 2: Advanced Customization via CSS

For complete control over every element in the interface—including fonts, spacing, hiding elements, or fine-tuning Dark Mode—you can provide a `custom_css_url` in your configuration.

This CSS file is loaded after all default stylesheets, giving your rules priority.

###### How to Identify Classes

Since the application is actively developed, class names and structures may evolve. We recommend using your browser's developer tools to identify the correct selectors for the elements you wish to modify.

1.  **Right-click** on the element you want to customize.
2.  Select **Inspect** (or "Inspect Element").
3.  Identify the class name or ID in the Elements panel.
4.  Write your CSS rule targeting that selector.

**Example:**

If you wanted to change the font size of the header title, you would inspect the title element, find its class (e.g., `.header-title`), and add the following to your custom CSS file:

```css
.header-title {
    font-size: 24px !important;
}
```

> **Note:** Using `!important` may be necessary to override certain default styles, especially utility classes.

<a id="page-28"></a>

## مسیر توسعهٔ اتصال اختصاصی

توسعه‌دهنده می‌تواند ایجاد اتاق، ورود کاربران، برند، ضبط و artifact را از برنامهٔ خود کنترل کند. مسیر پایه از مطالعهٔ احراز API، SDK و نمونهٔ PHP شروع می‌شود. در آی‌نشست، Core مسئول تماس مدیریتی و صدور توکن بعد از کنترل مجوز است و پنل انتخاب ارائه‌دهنده و نمایش خروجی را در اختیار کاربر می‌گذارد.

این بخش صفحهٔ `developer-guide/intro` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/developer-guide/intro) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/developer-guide/intro.md) · [نسخهٔ محلی کامل](docs/upstream/docs/developer-guide/intro.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Developer Guide Introduction

Welcome to the plugNmeet Developer Guide. This section is designed for developers and system integrators who want to move beyond the standard plugins and deeply integrate plugNmeet into their own applications.

Whether you want to build a fully white-labeled video conferencing solution, create complex business logic, or contribute to the open-source project, this guide is your starting point.

##### What You Can Achieve

plugNmeet is more than just a video conferencing tool; it's a flexible WebRTC platform built on a modern, API-driven architecture. With it, you can:

*   **Build a Custom, White-Label Experience:** Use our `getClientFiles` API to embed and completely re-brand the plugNmeet front-end client directly within your own application, without being trapped in an `iframe`.
*   **Automate and Manage Rooms:** Programmatically create, manage, and authenticate users for meetings using our comprehensive server-side APIs.
*   **Extend Functionality with Webhooks:** Listen for real-time events within a meeting (such as user join/leave, recording start/stop, etc.) to trigger custom workflows in your application.
*   **Scale for Massive Audiences:** Learn the architectural concepts required to build a distributed, high-availability deployment capable of supporting a massive user base.

##### How to Get Started

1.  **Explore the API Introduction:** This is the best place to start. It provides a complete reference for all available server endpoints and their parameters.
2.  **Review the Design Customization Guide:** Learn how to use a simple JavaScript object to completely change the look and feel of the front-end client.
3.  **Check out our SDKs:** For faster integration, we provide official SDKs for popular languages like PHP and JavaScript.

We're excited to see what you build with plugNmeet!

<a id="page-29"></a>

## اتصال اپلیکیشن موبایل و حالت Hybrid

در حالت hybrid رابط وب در WebView باقی می‌ماند و انتشار رسانه به لایهٔ native واگذار می‌شود. توکن ورود با client_type مناسب صادر می‌شود و bridge دوطرفه وضعیت نشست و فرمان‌ها را منتقل می‌کند. هویت مشترک میان native و وب، تحویل کلید E2EE، مجوزهای دستگاه و چرخهٔ خروج باید مطابق قرارداد این صفحه اجرا شوند. نام پیام‌های bridge در پیوست کامل آمده است.

این بخش صفحهٔ `developer-guide/mobile-app-integration` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/developer-guide/mobile-app-integration) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/developer-guide/mobile-app-integration.md) · [نسخهٔ محلی کامل](docs/upstream/docs/developer-guide/mobile-app-integration.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Build a Custom Video Conference App with plugNmeet Hybrid Integration

plugNmeet's hybrid integration model provides a fast and flexible path for building custom video conferencing apps for native platforms like iOS, Android, Flutter, React Native, or desktop.

This architecture combines a native media layer with the complete plugNmeet web client, which runs inside a WebView. The native application handles all media publishing (microphone, camera, and screen share), while the web client provides the feature-rich meeting interface. This approach allows developers to leverage native performance for media handling while avoiding the need to build a complex conferencing UI from scratch. It is especially useful for enabling features like screen sharing on mobile, which is not well-supported in mobile browsers.

##### See it in Action: Demo Applications

Before diving into the technical details, you can see this hybrid integration in action by exploring our open-source demo applications.

You might notice that plugNmeet does not provide its own mobile SDK. This is a deliberate design choice. Because plugNmeet uses LiveKit for its media layer, you can directly use LiveKit's excellent native SDKs (for iOS, Android, React Native, etc.) to handle all media-related tasks.

Instead of a restrictive SDK, we provide a flexible **pattern**: run the feature-rich plugNmeet web client in a WebView for the UI, and use LiveKit's native SDKs for the media layer. Our demo apps show just how simple this is to implement, giving you native performance and full control without having to build a UI from scratch.

Check the demo apps in our [plugNmeet mobile app repository](https://github.com/mynaparrot/plugnmeet-mobile-app) to understand the pattern.

##### Architecture and Identity

The core of the hybrid model is representing a single logical user with two participants in the LiveKit room:

| Identity | Used By | Permissions |
| --- | --- | --- |
| `[userID]` | Web client in WebView | Subscribe-only |
| `[userID]-native` | Native app | Publish-only |

The web client renders the UI and subscribes to all media streams, while the native app connects separately to publish its own media. The web client recognizes that tracks from a `*-native` identity matching its own `userID` should be treated as local media, displaying them in the correct UI tile without echo. This allows the existing web UI to function with minimal changes.

![Hybrid Integration Architecture](https://www.plugnmeet.org/img/hybrid-architecture.png)

##### End-to-End Workflow

###### Step 1: Get an Access Token

The workflow starts on your backend. To authorize a user, your server calls the plugNmeet API to generate a join token. For a hybrid app, it is crucial to set the `client_type` in this request.

1.  **Backend:** Call [`POST /room/getJoinToken`](https://www.plugnmeet.org/docs/api/room/join).
2.  **In the request body:** Set `user_info.client_type = "HYBRID_WEB"`.
3.  **Backend:** Receive the JWT (`access_token`) from the plugNmeet server and pass it securely to your native app, along with your `serverUrl`.

###### Step 2: Load the Web Client in a WebView

Once the native app has the `access_token`, it must load the plugNmeet web client into a WebView. There are two primary methods for this:

###### Option A: Load the Hosted Web Client Directly (Simple)

This is the easiest approach. Construct the URL with the access token and load it.

```kotlin
// Example for Android
val url = "${serverUrl.trimEnd('/')}/?access_token=${jwt}"
webView.loadUrl(url)
```

To apply a custom look and feel, you can add a `custom_design` query parameter to the URL. The value should be a URL-encoded JSON string. This is often used to change the logo and colors.

For a full list of supported parameters, please refer to our [Design Customisation Guide](https://www.plugnmeet.org/docs/developer-guide/design-customisation) for details.

```kotlin
// Example for Android with custom design
// The custom_design value is a URL-encoded JSON string like:
// {"logo":"https://your-domain.com/logo.png","primary_color":"#004D90"}
val customDesign = "%7B%22logo%22%3A%22https%3A%2F%2Fyour-domain.com%2Flogo.png%22%2C%22primary_color%22%3A%22%23004D90%22%7D"

val url = "${serverUrl.trimEnd('/')}/?access_token=${jwt}&custom_design=${customDesign}"
webView.loadUrl(url)
```

###### Option B: Build a Custom HTML Shell (Advanced)

This gives you maximum control. The native app calls `POST /api/getClientFiles` using the JWT for authentication (e.g., `Authorization: <access_token>`). This endpoint is designed for native apps and provides the same assets as the server-side `/auth/getClientFiles` endpoint, which uses API Key/Secret authentication. See the [Get Client Files API](https://www.plugnmeet.org/docs/api/get-client-files) for response details. The app then builds an HTML shell in memory from the returned assets and loads it into the WebView.

Even with this method, client-side customizations can be applied by injecting a `window.plugNmeetConfig` object. See our [Design Customisation Guide](https://www.plugnmeet.org/docs/developer-guide/design-customisation) for details.

###### Step 3: Enter Hybrid Mode

The web client validates the session by calling `POST /api/verifyToken`. It then enters hybrid mode based on one of two conditions:

-   **Token-Based (Recommended):** The response from `verifyToken` contains `client_type = HYBRID_WEB`.
-   **Configuration Override:** The native app injects `force_hybrid_web: true` into the `window.plugNmeetConfig` object. This forces hybrid mode regardless of the token's `client_type` and is useful for testing or custom URL schemes.

###### Step 4: Hand-Off to Native Publisher

In hybrid mode, the plugNmeet server provides two LiveKit tokens to the web client: a subscribe-only token for itself and a publish-only token for the native twin.

The web client connects to LiveKit, then passes the publish-only token to the native app over the communication bridge using the `INITIALIZE_NATIVE_PUBLISHER` action. Your native app now has everything it needs to connect to LiveKit and start publishing media.

##### Communication Bridge API

Communication between the WebView and the native host is handled by a message bridge defined in [`plugnmeet_native_bridge.proto`](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_native_bridge.proto). The wire format is **proto3 JSON text**.

###### Web Client to Native Host

| Action | Payload | Description |
| --- | --- | --- |
| `INITIALIZE_NATIVE_PUBLISHER` | `{ livekitUrl, token, nativeUserId, e2ee? }` | Provides the native host with the LiveKit URL and publish-only token. |
| `PUBLISH_NATIVE_MEDIA` | `{ source }` | Requests native publishing for `MIC`, `WEBCAM`, or `SCREENSHARE`. |
| `UNPUBLISH_NATIVE_MEDIA` | `{ source }` | Requests that the native app stop publishing the selected media source. |
| `MUTE_NATIVE_MEDIA` | `{ source }` | Mutes an already-published native track without unpublishing it. |
| `UNMUTE_NATIVE_MEDIA` | `{ source }` | Unmutes a previously muted native track. |
| `NATIVE_HEARTBEAT_PING` | `{ ts }` | Web-to-native keepalive message. |
| `TEARDOWN_NATIVE_PUBLISHER` | `{}` | Requests that the native app disconnect from LiveKit and release media resources. |

###### Native Host to Web Client

| Action | Payload | Description |
| --- | --- | --- |
| `NATIVE_MEDIA_STATUS` | `{ source?, deviceId?, error? }` | Confirms a native action or reports an error. |
| `NATIVE_TRACK_PUBLISHED` | `{ userId, kind, source }` | Notifies the web client that a native track was published. |
| `NATIVE_TRACK_UNPUBLISHED` | `{ userId, kind, source }` | Notifies the web client that a native track was unpublished. |
| `NATIVE_MEDIA_MUTED` | `{ source, muted }` | Confirms the native mute state. |
| `NATIVE_HEARTBEAT_PONG` | `{ ts }` | Native response to heartbeat ping. |
| `NATIVE_ERROR` | `{ msg, context? }` | Reports a generic error from the native host. |

##### E2EE Hand-Off

If end-to-end encryption is enabled, the plain text encryption key is handed to the native app as part of the `INITIALIZE_NATIVE_PUBLISHER` payload.

```json
{
  "action": "INITIALIZE_NATIVE_PUBLISHER",
  "initializeNativePublisher": {
    "livekitUrl": "wss://livekit.example.com",
    "token": "<native_publish_only_livekit_token>",
    "nativeUserId": "<userID>-native",
    "e2ee": {
      "enabled": true,
      "key": "your-e2ee-plain-text-key"
    }
  }
}
```

##### Native App Implementation Checklist

###### Required

-   Receive `serverUrl` and plugNmeet JWT from your backend.
-   Load the plugNmeet web client in a WebView.
-   Implement the bridge logic to send and receive messages.
-   Handle `INITIALIZE_NATIVE_PUBLISHER` by connecting to LiveKit with the provided token.
-   Implement native media capture and publishing for `PUBLISH_NATIVE_MEDIA` (mic & camera).
-   Implement `MUTE_NATIVE_MEDIA`, `UNMUTE_NATIVE_MEDIA`, and `UNPUBLISH_NATIVE_MEDIA`.
-   Implement `TEARDOWN_NATIVE_PUBLISHER` to clean up resources.
-   Send status messages back to the web client (`NATIVE_MEDIA_STATUS`, `NATIVE_TRACK_PUBLISHED`, etc.).
-   Reply to `NATIVE_HEARTBEAT_PING` with `NATIVE_HEARTBEAT_PONG`.
-   Implement a watchdog to tear down the native publisher if the WebView heartbeat stops.

###### Recommended

-   Implement `PUBLISH_NATIVE_MEDIA` for screen sharing.
-   Treat repeated `INITIALIZE_NATIVE_PUBLISHER` calls as a full reinitialization.
-   Provide clear native permission prompts and error handling.
-   Avoid embedding API keys or secrets in production apps.

###### Optional

-   E2EE support.
-   Native virtual background or camera effects.
-   A custom device picker UI.
-   Background mode and audio routing controls.

<a id="page-30"></a>

## مقیاس‌پذیری و جداسازی ضبط‌کننده

پیشنهاد صفحه، بررسی انتقال ضبط‌کننده به ماشین مستقل به عنوان نخستین گام رایج رشد است. پس از آن ارتقای عمودی و سپس توزیع لایه‌های برنامه و رسانه مطرح می‌شود. NATS، دیتابیس و Redis به پایداری و تنظیمات خوشه‌ای متناسب نیاز دارند. سرور و recorder باید به فایل‌ها در مسیر سازگار دسترسی داشته باشند و فضای مشترک یا hook ذخیره‌سازی، دسترس‌پذیری خروجی را تأمین کند.

این بخش صفحهٔ `developer-guide/scalable-setup` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/developer-guide/scalable-setup) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/developer-guide/scalable-setup.md) · [نسخهٔ محلی کامل](docs/upstream/docs/developer-guide/scalable-setup.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Scaling Your PlugNmeet Deployment

This guide provides a high-level overview and recommended strategies for scaling a self-hosted plugNmeet deployment to support a large number of concurrent users across multiple servers.

###### Before You Begin: A Phased Approach to Scaling

It is critical to understand that a fully distributed, multi-server setup is the most advanced and complex configuration. For many users, simpler scaling strategies can provide the necessary performance boost without the added complexity.

We recommend approaching scaling in phases:

###### Phase 1: Isolate the Recorder (The Most Common Bottleneck)

For optimal performance, especially when recording frequently, we strongly recommend deploying the **`plugnmeet-recorder` on a dedicated server**, separate from your main `plugnmeet-server` and LiveKit instance.

*   **Why?** Recording and post-processing (transcoding) are CPU-intensive tasks. Isolating them prevents recording jobs from impacting the performance of your live meetings.
*   **When?** This should be your first step if you notice performance degradation during active recordings.
*   **Key Requirements:** When you isolate the recorder, you must configure a **shared network storage** for your recording files and assign a **unique ID** to the recorder instance.
*   **Next Steps:** For detailed configuration instructions, please refer to [The plugNmeet Recorder](https://www.plugnmeet.org/docs/developer-guide/scalable-setup#3-the-plugnmeet-recorder) section later in this guide.

###### Phase 2: Vertically Scale Your Main Server

If you are not recording heavily but are still seeing performance issues, the next step is to **vertically scale** your single server. This simply means adding more CPU cores and RAM. A single, well-provisioned server can comfortably handle **hundreds of concurrent users** — see our measured [Benchmarks](https://www.plugnmeet.org/docs/benchmark) (300–500 in one room clean on 4 vCore / 16 GB, ~900 limit).

###### Phase 3: Advanced Distributed Setups

**For Advanced, Large-Scale Deployments**

Most users will **never** need to go beyond Phase 1 (isolating the recorder) and Phase 2 (vertically scaling). A single, well-provisioned server can handle a very large number of concurrent users and operations.

The fully distributed architecture described in the rest of this guide is designed for massive, enterprise-scale deployments with extreme high-availability requirements. You should only consider this phase after you have exhausted vertical scaling options.

When you do reach this scale, **remember that it is not an all-or-nothing process**. You can incrementally scale individual components (like LiveKit or NATS) as needed, rather than clustering everything at once. The rest of this guide details the conceptual architecture for when you need to scale each part.



---

##### Conceptual Architecture

In a distributed setup, the core components of plugNmeet are separated onto different machines. A typical large-scale architecture consists of:

1.  **Stateless Application & Media Layers:**
    *   Multiple `plugnmeet-server` instances.
    *   Multiple `LiveKit` media server instances.
    *   A standard Load Balancer.
2.  **Stateful & Clustered Infrastructure Layers:**
    *   A `NATS` cluster.
    *   A `Redis` cluster.
    *   A `MariaDB` cluster.
3.  **Horizontally Scaled Recorders:**
    *   Multiple `plugnmeet-recorder` instances, potentially in different operational modes.
4.  **Shared & Persistent File Storage:**
    *   A shared network filesystem (e.g., NFS, Samba (CIFS), GlusterFS, or S3-compatible) accessible by all `plugnmeet-server` and `plugnmeet-recorder` instances. This is a critical component that ensures recordings, analytics files, and user uploads are consistently available across the entire cluster.

---

##### Component Scaling Strategy

The key to a successful distributed setup is understanding how each component scales. The components fall into two main categories: stateless services that can be easily load-balanced, and stateful services that require their own specific clustering configurations.

###### 1. Stateless Application & Media Layers

These components can be placed behind a standard load balancer, whose primary role is to terminate SSL and distribute initial traffic.

###### The plugNmeet Server: The Lightweight Command Post

The main `plugnmeet-server` application is designed to be a lightweight, stateless **command post** or orchestrator. It intentionally offloads all heavy lifting: media is handled by LiveKit, and recording is handled by the `plugnmeet-recorder`.

Its core responsibility is to handle API requests, authenticate users, and coordinate the overall session via NATS. Because it uses very few resources (CPU and RAM), it is incredibly efficient and easy to scale horizontally.

*   **Action:** You can run multiple instances on different machines and place them behind a standard L4 or L7 load balancer using a simple algorithm like round-robin or least connections.
*   **Configuration:** Each instance must be configured to point to the same shared infrastructure (NATS, Redis, and MariaDB).

###### LiveKit Media Server

LiveKit also works behind a load balancer for initial connections.
*   **Action:** The load balancer can terminate SSL and proxy the initial connection to any available LiveKit node.
*   **Critical Note:** Once the connection is established, LiveKit's own internal clustering logic handles all complex media routing. It is **essential to follow the official LiveKit documentation for clustering** to configure the backend correctly.
*   **Firewall Configuration:** Ensure that your firewall rules allow the necessary TCP and UDP ports for WebRTC traffic to reach your LiveKit instances. For a detailed list of required ports, please refer to our **[Firewall Configuration Guide](https://www.plugnmeet.org/docs/firewall)**.
*   **Dependency:** A standalone or clustered Redis instance is required for LiveKit clustering.

**(Link: [Official LiveKit Clustering Guide](https://docs.livekit.io/home/self-hosting/distributed))

---

###### 2. Stateful & Clustered Infrastructure Layers

**Critical Warning:** The following components have their own robust clustering and service discovery mechanisms. Placing them behind a standard, simple TCP/L4 load balancer without specific, advanced configuration is **highly likely to cause system instability, data corruption, or complete failure.** Always follow the official clustering documentation for these services.

###### NATS Messaging System

*   **Recommendation:** NATS has its own clustering and client routing capabilities. Your `plugnmeet-server` and `plugnmeet-recorder` clients should be configured with a list of all NATS nodes in the cluster. The clients will automatically handle failover.
*   **Action:** Follow the **official NATS documentation for creating a cluster**.

**(Link: [Official NATS Clustering Documentation](https://docs.nats.io/running-a-nats-service/configuration/clustering/jetstream_clustering))

###### MariaDB Database

*   **Recommendation:** For high availability, the database must be clustered. All `plugnmeet-server` instances must point to the same database cluster.
*   **Action:** Use a managed, high-availability database service from a cloud provider, or set up your own solution like a MariaDB Galera Cluster. Follow the official documentation for your chosen solution.

**(Link: [Official MariaDB Documentation](https://mariadb.com/docs/galera-cluster))

###### Redis Cache

*   **Recommendation:** For a highly available plugNmeet and LiveKit cluster, your Redis instance must also be highly available.
*   **Action:** Use a managed Redis service or set up your own Redis Cluster or Sentinel configuration. Your plugNmeet and LiveKit nodes should be configured to connect to this cluster directly.

**(Link: [Official Redis Sentinel Clustering Documentation](https://redis.io/docs/latest/operate/oss_and_stack/management/sentinel))

---

###### 3. The plugNmeet Recorder

The `plugnmeet-recorder` is a special case. While it is stateful during a job, its discovery and load balancing are handled automatically.

*   **Load Balancing & Discovery:** The recorder does **not** require an external load balancer. Recorders automatically announce their availability over NATS, and the `plugnmeet-server` distributes jobs to them.

*   **Operational Modes:** The `plugnmeet-recorder` supports different operational modes, allowing for a highly scalable and resilient pipeline. Each instance can be configured via the `mode` setting in `config.yaml`:
    *   **`recorderOnly`:** This instance will *only* handle live session recording. Once a raw recording file is captured, it publishes a transcoding job to a queue and immediately becomes available for the next live session. This is ideal for your primary recording servers as it keeps them lightweight and responsive.
    *   **`transcoderOnly`:** This instance will *only* process transcoding jobs. It subscribes to the job queue, performs the CPU-intensive conversion of raw files to MP4, and does not handle any live recordings. You can run a fleet of these workers on cheaper, CPU-optimized VMs to process recordings in parallel without impacting live meetings.
    *   **`both` (Default):** A single instance performs both live recording and transcoding. This is suitable for smaller setups or if you record infrequently.

*   **Configuration:**
    *   Each recorder instance must be configured to connect to the same NATS cluster.
    *   Each recorder instance **must** be assigned a unique `id` in its `config.yaml` under the `recorder` section (e.g., `recorder.id: "recorder-node-1"`). This ID is used by the `plugnmeet-server` to track and communicate with individual workers. It is especially critical for `recorderOnly` and `both` modes, and highly recommended for all instances for clear logging and identification.

**Administrator Responsibility**

It is **your responsibility** to ensure that every single recorder instance has a globally unique `id`. The `plugnmeet-server` does **not** validate the uniqueness of these IDs.

If multiple recorders share the same ID, they will all try to accept the same recording job from the NATS queue. This race condition will cause the job to fail, as none of the instances will be able to acquire the necessary lock. This scenario is extremely difficult to debug, as the only symptom may be that recordings intermittently fail to start.



*   **Action:** For detailed setup, refer to the **official `plugnmeet-recorder` repository README**.

**(Link: [Official plugNmeet-recorder README](https://github.com/mynaparrot/plugNmeet-recorder))

---

###### 4. Shared Filesystem: A Critical Requirement for Data Consistency

In a distributed environment, it is **absolutely critical** that all stateless servers have a consistent way to access persistent files. Without this, your users will experience broken features like missing file downloads or inaccessible recordings.

There are two primary approaches to achieve this: a traditional shared filesystem or the more modern and flexible hooks system.

###### Option 1: Traditional Shared Filesystem

You can mount a shared network storage solution (such as **NFS**, **Samba (CIFS)**, **GlusterFS**, or an **S3-compatible object store** mounted as a filesystem) to the same path on **all** of your `plugnmeet-server` and `plugnmeet-recorder` instances.

The following paths in your `config.yaml` must point to this shared location:

*   `recorder_info.recording_files_path`
*   `artifacts_settings.storage_path`
*   `upload_file_settings.path`

###### Option 2: The Hooks System (Recommended for Cloud-Native Deployments)

**Alternative to Shared Filesystems**
For a more flexible and cloud-native approach, you can use the **[Scripting & Storage Hooks](https://www.plugnmeet.org/docs/others/hooks)** system instead of a shared filesystem. This is the recommended method for handling recordings and artifacts in a scalable environment.

With hooks, you can directly integrate with object storage services like **AWS S3**, **Google Cloud Storage**, or **MinIO**.

*   **For Recordings:** Your `post_transcoding` hook on the recorder uploads the final MP4 to your object store and returns a unique identifier (like an S3 key).
*   **For Artifacts & Downloads:** Your `upload_hook` on the server uploads artifacts to the same object store. Your `download_hook` and `delete_hook` on the server then use the stored identifier to generate secure download links or delete files directly from the object store.

This approach eliminates the need for managing a shared network disk, which is often more complex and less scalable than using a dedicated object storage service.


---

##### Configuration

Once your distributed infrastructure is in place, you simply need to update your `config.yaml` file for each `plugnmeet-server` instance to point to the new clustered services.

```yaml
# Example config snippet
redis_info:
  host: "your-redis-cluster-endpoint:6379"
  # ... other redis settings

database_info:
  host: "your-mariadb-cluster-endpoint"
  port: 3306
  username: "your_user"
  password: "your_password"
  db: "plugnmeet"

nats_info:
  nats_urls:
    - "nats://node1.your-domain.com:4222"
    - "nats://node2.your-domain.com:4222"
    - "nats://node3.your-domain.com:4222"

livekit_info:
  host: "http://your-livekit-endpoint" # The endpoint for your LiveKit cluster
  # ... other livekit settings
```

By following this conceptual guide and leveraging the official documentation for each component, you can build a robust, scalable plugNmeet deployment capable of supporting a massive user base. For measured capacity at each stage — laptop, single server, cluster — see [Benchmarks](https://www.plugnmeet.org/docs/benchmark).

<a id="page-31"></a>

## محیط توسعهٔ محلی

سه مخزن سرور، کلاینت و ضبط‌کننده برای توسعهٔ کامل استفاده می‌شوند. نمونهٔ پیکربندی پیش از اجرا کپی و متناسب با محیط اصلاح می‌شود. سرویس‌های زیرساختی با Docker اجرا می‌شوند و کلاینت با ابزارهای Node و pnpm توسعه می‌یابد. نسخهٔ Go و Node باید از فایل‌های همان commit خوانده شود. URLهای قابل دسترس مرورگر و URLهای داخلی کانتینر ممکن است متفاوت باشند و اتصال LiveKit و NATS هر دو باید برقرار باشد.

این بخش صفحهٔ `developer-guide/setup-development` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/developer-guide/setup-development) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/developer-guide/setup-development.md) · [نسخهٔ محلی کامل](docs/upstream/docs/developer-guide/setup-development.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Setting Up Your Development Environment

Welcome to the plugNmeet development guide! We're excited that you're interested in contributing. This guide will walk you through setting up a complete local development environment, allowing you to run the entire plugNmeet stack on your own machine.

Following these steps will enable you to test new features, fix bugs, and see your changes in real-time.

---

##### Prerequisites

Before you begin, please ensure you have the following software installed on your computer:

1.  **Docker**: Used to run the backend services (like LiveKit and NATS) in isolated containers, which simplifies the setup process significantly.
2.  **Node.js**: Required for running and developing the `plugNmeet-client`. We recommend using the latest LTS version.
3.  **Git**: For cloning the project repositories. You can use a command-line client or a GUI like GitHub Desktop.
4.  **FFMPEG**: (Optional) Required only if you plan to work on or test the recording features.

---

##### Step 1: Clone the Repositories

First, create a main project directory and clone the necessary plugNmeet repositories into it. Open your terminal and run the following commands:

```sh
# Create and enter a new directory for the project
mkdir plugnmeet-dev
cd plugnmeet-dev

# Clone the main server repository
git clone https://github.com/mynaparrot/plugNmeet-server server

# Clone the front-end client repository
git clone https://github.com/mynaparrot/plugNmeet-client client

# Clone the recorder service repository
git clone https://github.com/mynaparrot/plugNmeet-recorder recorder
```

---

##### Step 2: Configure and Run the Backend Server

The `plugNmeet-server` repository contains a `docker-compose.yaml` file that orchestrates all the necessary backend services.

```sh
# Navigate to the server directory
cd server

# For the first time setup, copy the sample configuration files
cp config_sample.yaml config.yaml
cp livekit_sample.yaml livekit.yaml
cp ingress_sample.yaml ingress.yaml
cp nats_server_sample.conf nats_server.conf
cp docker-compose_sample.yaml docker-compose.yaml

# Build and start all backend services
docker-compose up --build
```

The first time you run this, it may take a few minutes to download and build the container images. On subsequent runs, you can simply use `docker-compose up` to start the services much faster.

Leave this terminal window running, as it will show the logs for all backend services.

---

##### Step 3: Configure and Run the Front-End Client

Now, let's get the user interface running. Open a **new terminal window or tab** so you can run the client simultaneously with the server.

```sh
# Navigate to the client directory from your main project folder
cd client

# Copy the sample configuration file
cp src/assets/config_sample.js src/assets/config.js

# Install all project dependencies
# If you don't have pnpm, you can install it with `npm install -g pnpm`
pnpm install

# Start the client development server
pnpm start
```

Once the development server is running, you can access the plugNmeet client in your web browser at: **[http://localhost:3000/login.html](http://localhost:3000/login.html)**

---

##### Step 4: (Optional) Run the Recorder Service

If you need to work on the recording features, you can start the recorder service. Open a **third terminal window or tab**.

```sh
# Navigate to the recorder directory from your main project folder
cd recorder

# Copy the sample configuration files
cp config_sample.yaml config.yaml
cp docker-compose_sample.yaml docker-compose.yaml

# Build and start the recorder service
docker-compose up --build
```

The recorder will now be running and will automatically listen for recording jobs from the `plugNmeet-server`.

---

##### Troubleshooting

###### "I can't connect to the meeting room."

This is a common issue on local development setups, often caused by Docker's networking. The LiveKit server inside the container needs to know your host machine's local IP address to establish a WebRTC connection.

1.  Find your computer's local IP address (e.g., `192.168.1.100`).
2.  Open the `docker-compose.yaml` file inside the `server` directory.
3.  Find the `livekit` service section.
4.  Modify the `command` line to add the `--node-ip` flag with your IP:

```yaml
  # Before
  command: --config "/app/livekit.yaml" --dev

  # After
  command: --config "/app/livekit.yaml" --dev --node-ip YOUR_LOCAL_IP_HERE
```

5.  Restart the server with `docker-compose down` and `docker-compose up`.

---

##### What's Next?

You now have a fully functional plugNmeet development environment! You can start making changes to the client or server code, and your changes will be reflected in real-time.

We welcome all contributions. A great place to start is by looking at the open issues on our GitHub repositories. Happy coding!

<a id="page-32"></a>

## فایروال و پورت‌های رسانه

ورودی عمومی HTTP روی 80/TCP برای صدور گواهی و HTTPS روی 443/TCP برای برنامه و API لازم است. رسانه از 50000 تا 60000/UDP و مسیر جایگزین 7881/TCP عبور می‌کند. TURN نصب رسمی از 443/TCP و 443/UDP استفاده می‌کند. در میزبان مشترک، مالکیت این پورت‌ها و مسیریابی دامنه باید پیش از افزودن سرویس دیگر ثبت شود. پورت‌های SIP فقط در صورت فعال‌سازی تماس تلفنی کاربرد دارند. آزمون واقعی مرورگر و TURN، عبور رسانه را تأیید می‌کند.

این بخش صفحهٔ `firewall` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/firewall) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/firewall.md) · [نسخهٔ محلی کامل](docs/upstream/docs/firewall.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Firewall & Ports Configuration

A correctly configured firewall is the most critical step for a successful Plug-N-Meet installation. If your server is behind a firewall or NAT (which is true for almost all cloud providers like AWS, Google Cloud, Azure, etc.), you must ensure that the correct ports are open to allow traffic to reach the server.

This guide will show you which ports to open and how to test your configuration.

---

##### Required Ports

You must allow incoming traffic on the following ports:

| Port          | Protocol | Required | Description                                                  |
| :------------ | :------- | :------- | :----------------------------------------------------------- |
| **80**        | TCP      | **Yes**  | Required by Let's Encrypt to issue SSL certificates.         |
| **443**       | TCP      | **Yes**  | The primary HTTPS port for all application and API traffic.  |
| **7881**      | TCP      | **Yes**  | A fallback port for WebRTC media when UDP is blocked.        |
| **50000-60000** | UDP      | **Yes**  | The primary port range for all WebRTC audio and video media. |

**What about TURN?**
The installation script automatically sets up and configures a TURN server for you on port 443 over TCP/UDP, so you typically do not need to open a separate port for it.


---

###### SIP/VoIP Dial-In Ports

If you plan to use the SIP/VoIP Dial-In feature, you must also open the following ports to allow communication with your external SIP trunking provider.

| Port          | Protocol | Required | Description                                                  |
| :------------ | :------- | :------- | :----------------------------------------------------------- |
| **5060-5061** | TCP      | **Yes**  | SIP signaling ports.                                         |
| **5060-5061** | UDP      | **Yes**  | SIP signaling ports.                                         |
| **10000-20000** | UDP      | **Yes**  | RTP (Real-time Transport Protocol) media port range for audio. |

**Further Details**
For a deeper understanding of the SIP gateway's networking requirements and advanced configurations, please consult the official [livekit/sip GitHub repository](https://github.com/livekit/sip).


---

##### Cloud Provider Guides

Most cloud providers have a "Security Group" or "Firewall" section where you can create these rules. Here are links to the official documentation for popular providers:

*   [Amazon AWS](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/authorizing-access-to-an-instance.html)
*   [Google Cloud](https://cloud.google.com/vpc/docs/using-firewalls)
*   [Microsoft Azure](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/nsg-quickstart-portal)
*   [DigitalOcean](https://docs.digitalocean.com/products/networking/firewalls/how-to/configure-rules/)

---

##### Testing Your Configuration (Optional but Recommended)

Before running the main installation script, it's a good idea to verify that your ports are open. The `ncat` tool (an improved version of `netcat`) is perfect for this.

###### 1. On Your Plug-N-Meet Server

First, install `ncat` and start listening on a port. We'll use port 443 as an example.

```sh
# Install the tool
sudo apt update && sudo apt install -y nmap

# Listen for TCP traffic on port 443
sudo ncat -l -p 443 --keep-open --verbose
```

The server will now wait for a connection.

###### 2. From an External Computer

On a *different computer* (like your laptop), use `ncat` to connect to your server. Replace `your_server_ip` with your server's public IP address.

```bash
# On macOS (with Homebrew) or Linux: brew install nmap
# On Windows: download from nmap.org

ncat your_server_ip 443 --verbose
```

Now, type `hello` and press Enter.

###### 3. Check the Results

*   If the connection is **successful**, you will see "Connection from [your_ip] received" on your server, and the word `hello` will appear.
*   If the connection **fails** (it times out or is refused), your firewall rule is incorrect. Double-check your cloud provider's settings and try again.

You can repeat this test for the other ports (`80`, `7881`) and for UDP by adding the `-u` flag (e.g., `sudo ncat -u -l -p 50005 ...`).

<a id="page-33"></a>

## نصب و نگهداری پایه

نصب‌کنندهٔ رسمی برای میزبان تمیز Ubuntu یا Debian نوشته شده است. ورودی‌های لازم شامل IP عمومی، دامنهٔ جلسه، دامنهٔ TURN و تنظیمات گواهی هستند. نصب‌کننده پیکربندی چند سرویس را تغییر می‌دهد، بنابراین روی میزبان مشترک ابتدا محتوای آن بازبینی و تنظیمات به Compose اختصاصی منتقل می‌شود. حداقل منابع تولیدی پیشنهادی صفحه، ۴ هسته و ۴ گیگابایت RAM است و برای ضبط ۸ هسته و ۸ گیگابایت یا بیشتر توصیه شده است. ظرفیت عملی با بار، پهنای باند و تعداد ضبط هم‌زمان تعیین می‌شود.

این بخش صفحهٔ `installation` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/installation) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/installation.md) · [نسخهٔ محلی کامل](docs/upstream/docs/installation.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Installation Quick Start

This guide will walk you through the fastest and easiest way to get your own self-hosted video conferencing platform up and running.

Our recommended method uses a simple, automated script that sets up all the necessary components required to run the platform on a clean Ubuntu or Debian server. Alternatively, for those who prefer a managed service, our official **[PlugNmeet Cloud solution](https://www.plugnmeet.cloud)** provides a ready-to-use server, allowing you to bypass the installation process entirely.

Let's get your server set up.

---

##### Step 1: Server Requirements

Before you begin, make sure you have the following:

*   A **clean Ubuntu or Debian server** with a **public IP address**. We recommend any Ubuntu LTS release. It is critical that the server does **not** have other web servers like Apache or Nginx already running, as this will cause conflicts.
*   **Hardware & Bandwidth:** The Plug-N-Meet server itself is a very lightweight application. The hardware you need depends on your expected workload, specifically the number of concurrent users and your recording frequency.

    For a small group or testing, a simple server (e.g., 2 CPU cores, 2-4GB RAM) is sufficient.

    For a production environment, we recommend:
    *   **Bandwidth:** This is the most critical factor. A **minimum of 100 Mbits/sec** is recommended, but more is always better for a quality user experience.
    *   **CPU:** At least 4 cores (8+ recommended if using the recorder).
    *   **RAM:** At least 4GB (8GB+ recommended if using the recorder).

    Wondering what that handles? A single 4-vCore / 16 GB server ran 300–500 users in one room cleanly in our [Benchmarks](https://www.plugnmeet.org/docs/benchmark).
*   **Two subdomains** pointing to your server's public IP (e.g., `plugnmeet.example.com` and `turn.example.com`).
*   A **valid email address** for generating a Let's Encrypt SSL certificate.

**Firewall Configuration**
If your server is behind a firewall, please review our **[Firewall Configuration Guide](https://www.plugnmeet.org/docs/firewall)** to ensure the necessary ports are open before you proceed.


---

##### Step 2: Run the Automated Installation

Connect to your server via SSH. The recommended and most transparent method is to download the script first and then execute it. This allows you to inspect the script's contents before running it as root.

```bash
# 1. Download the script
wget https://raw.githubusercontent.com/mynaparrot/plugNmeet-install/main/install.sh
```

Then choose one of the two modes below.

###### Option A: Interactive install (easiest)

Just run the script with no arguments. It will prompt you for everything it needs, such as your domain names, email address, and whether to install the recorder and firewall.

```bash
# 2. Run the script as root
sudo bash install.sh
```

###### Option B: Non-interactive install (with options)

Pass everything as arguments and the script will run without prompting. This is ideal for automation or as a cloud post-install script (e.g. cloud-init user-data), so a fresh server can come up with plugNmeet fully ready.

```bash
sudo bash install.sh --domain plugnmeet.example.com --turn-domain turn.example.com --email admin@example.com --recorder yes --ufw yes
```

| Option | Description |
|---|---|
| `-d, --domain <domain>` | plugNmeet server domain (e.g. `plugnmeet.example.com`) |
| `-t, --turn-domain <domain>` | TURN server domain (e.g. `turn.example.com`) |
| `-e, --email <email>` | Email address for the Let's Encrypt SSL certificate |
| `-r, --recorder <yes\|no>` | Install the recorder (`yes` or `no`; if omitted, you will be prompted) |
| `-f, --ufw <yes\|no>` | Configure the UFW firewall (`yes` or `no`; if omitted, you will be prompted) |
| `-h, --help` | Show all options and exit |

A few things to know:

* Both `--option value` and `--option=value` forms work (e.g. `--domain=plugnmeet.example.com`).
* Omit any option and the script will simply prompt you for it.
* Run `sudo bash install.sh --help` at any time to see all options.
* If there is no terminal to prompt you (for example in a cloud post-install script) and a value is missing, the script stops with a clear error telling you which option to pass.

**Alternatively**, you can download and execute it in a single command:
```bash
sudo su -c "bash <(wget -qO- https://raw.githubusercontent.com/mynaparrot/plugNmeet-install/main/install.sh)" root
```

**Cloud post-install example**
Most cloud providers let you run a script automatically after the OS is installed (e.g. cloud-init user-data). Use the non-interactive form so the server comes up fully ready with no manual steps:

```bash
#!/bin/bash
wget -q https://raw.githubusercontent.com/mynaparrot/plugNmeet-install/main/install.sh -O /root/install.sh
bash /root/install.sh --domain plugnmeet.example.com --turn-domain turn.example.com --email admin@example.com --recorder yes --ufw yes
```


At the end of the process, the script will display your **Server URL**, **API Key**, and **API Secret**.

**Your Server Credentials**
We recommend saving these credentials now for easy access. However, if you happen to close the window, don't worry! You can always find them later in the `/opt/plugNmeet/config.yaml` file on your server.


---

##### Step 3: Verify Your Installation

Once the installation is complete, you can verify that everything is running correctly.

1.  **Check the Service Status:**

    ```bash
    # Check the main server
    systemctl status plugnmeet

    # Check the recorder service
    systemctl status plugnmeet-recorder
    ```
    Both services should show an `active (running)` status.

2.  **Perform a Test Login:**
    The best way to confirm that the entire system is working is to join a test meeting.

    Navigate to `https://plugnmeet.example.com/login.html`. On this page, you can use the `API_KEY` and `API_SECRET` you saved from the installation to generate a temporary access token and join a test conference.

---

##### What's Next?

Congratulations! You now have a fully functional, self-hosted video conferencing platform. Here are a few things you can do next:

*   **Start Using Plug-N-Meet:** Now that your server is running, visit our **[User Guide Overview](https://www.plugnmeet.org/docs/user-guide/overview)** to learn about the different ways you can integrate and use your new platform, from no-code plugins to our powerful API.
*   **Explore the API:** For developers ready to dive in, our **[API Documentation](https://www.plugnmeet.org/docs/api/intro)** provides all the details you need to start building.
*   **Manage Your Server:** Learn about the common commands for **[managing your services](https://www.plugnmeet.org/docs/installation#managing-services)**, including starting, stopping, and viewing logs.
*   **Advanced Customization**: Explore our **[Scripting & Storage Hooks](https://www.plugnmeet.org/docs/others/hooks)** to automate workflows and integrate with external storage providers.

---

##### Advanced Topics

###### Managing Services

All Plug-N-Meet components run as `systemd` services.

```bash
# To start, restart, or stop the services:
systemctl restart plugnmeet
systemctl restart plugnmeet-recorder
```

To view the log files:
```sh
tail -n 200 -f /opt/plugNmeet/log/plugNmeet.log
tail -n 200 -f /opt/plugNmeet/recorder/logs/recorder.log
```

###### Updating Your Server

To update your Plug-N-Meet installation to the latest version, simply run the `update.sh` script:

```sh
sudo su -c "bash <(wget -qO- https://raw.githubusercontent.com/mynaparrot/plugNmeet-install/main/update.sh)" root
```

###### Scalable Deployments

For very large-scale needs, Plug-N-Meet can be deployed in a distributed, multi-server cluster. To learn more, read our **[Scalable Deployment Guide](https://www.plugnmeet.org/docs/developer-guide/scalable-setup)**.

<a id="page-34"></a>

## معرفی و معماری ابزار

plugNmeet یک سامانهٔ خودمیزبان جلسهٔ آنلاین است که سرور Go، کلاینت React و ضبط‌کنندهٔ مستقل دارد. انتقال رسانه بر عهدهٔ LiveKit است و NATS پیام‌های بلادرنگ و همگام‌سازی جلسه را جابه‌جا می‌کند. قابلیت‌ها شامل صوت و تصویر، اشتراک صفحه، وایت‌برد، نظرسنجی، اتاق فرعی، ضبط و خروجی‌های تحلیلی است. اتصال سفارشی آی‌نشست از API و وب‌هوک انجام می‌شود و رابط جلسه می‌تواند روی دامنهٔ مستقل میزبانی شود.

این بخش صفحهٔ `intro` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/intro) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/intro.md) · [نسخهٔ محلی کامل](docs/upstream/docs/intro.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### plugNmeet: Self-Hosted Open-Source Video Conferencing for Your Website

Welcome to plugNmeet, the video conferencing platform designed for ultimate control, flexibility, and **cost-effective scale**. Host secure, high-quality, and **AI-powered** meetings directly from your website, with your own branding, and complete ownership of the user experience.

What truly sets plugNmeet apart is its **unmatched customization**. The interface is fully controllable, allowing you to tailor the entire experience — from simple branding with your own logo and colors, to enabling specific features to create anything from a minimalist video chat, to a full-featured **virtual classroom**, or even a massive live event with [**hundreds of participants**](https://www.plugnmeet.org/docs/benchmark).

**In 30 seconds**
- **What:** Complete meeting room — audio/video, screen share, whiteboard, chat, polls, breakout rooms, recording, live streaming, and AI summaries/translations.
- **How you own it:** Install the free server on your own server, or skip ops with [plugNmeet Cloud](https://www.plugnmeet.cloud). In both cases the meeting runs under your brand.
- **How you use it:** No-code plugins for WordPress, Moodle, Joomla, and any LMS via LTI — or plugNmeet API + SDKs for custom apps.
- **Why teams pick it:** Privacy by design with optional end-to-end encryption, true white-label via config (no code forks), and a lightweight, scalable stack.


##### Which one are you? Pick your path

| If you are… | You get… | Start here |
|---|---|---|
| **A creator, educator, or business owner** — you want video on your site without coding | Install a plugin, enter server URL + key, add meetings from your dashboard | [User Guide Overview](https://www.plugnmeet.org/docs/user-guide/overview) |
| **An admin or self-hoster** — you care about privacy, cost, and control | A single-binary Go server + automated install script, you own data and recordings | [Installation Quick Start](https://www.plugnmeet.org/docs/installation) |
| **A developer** — you want to embed, automate, or ship mobile apps | PlugNmeet API, PHP/JS SDKs, webhooks, custom mobile apps, and full UI theming without an iframe trap | [Developer Guide](https://www.plugnmeet.org/docs/developer-guide/intro) · [API Intro](https://www.plugnmeet.org/docs/api/intro) |

You can mix paths: many teams self-host once, then let teachers use the Moodle plugin while developers automate rooms via API.

##### What can you build with it?

Same engine, very different products — you just enable what you need:

- **Virtual classroom & coaching:** Whiteboard + shared notepad + polls + breakout rooms + attendance via [analytics](https://www.plugnmeet.org/docs/others/analytics). Works from Moodle, WordPress, or any LMS via [LTI](https://www.plugnmeet.org/docs/user-guide/lti).
- **Webinar & live event:** Waiting room, moderator controls, [RTMP live streaming](https://www.plugnmeet.org/blog/custom-qna-webinar-rtmp-broadcast-api) to YouTube/Facebook, [studio input from OBS](https://www.plugnmeet.org/blog/obs-rtmp-whip-ingress) via RTMP/WHIP, shared video playback for everyone.
- **Telehealth & private consulting:** Browser-only join, lockable room, per-user permissions, and optional [end-to-end encryption](https://www.plugnmeet.org/docs/security-overview) where even the server cannot read media or chat.
- **Community or marketplace calls:** Minimal video chat embedded in your product with your logo, colors, and domain. No plugNmeet branding.
- **Internal meetings & ERP integration:** API-created rooms, SSO via your existing users, recordings stored as plain MP4 on your storage.
- **Custom mobile & desktop apps:** Ship your own branded iOS, Android, Flutter, React Native, or desktop app with the [hybrid integration](https://www.plugnmeet.org/docs/developer-guide/mobile-app-integration) — native media performance plus the full plugNmeet web UI in a WebView, no mobile SDK lock-in.

If you are coming from BigBlueButton, see [Migrating from BigBlueButton](https://www.plugnmeet.org/docs/tutorials/migration-from-bbb) — plugNmeet speaks the BBB API, so existing apps keep working.

##### Key capabilities at a glance

You do not need to learn all of this on day one. This is the full toolbox your room can offer:

###### Meet without friction
- **HD audio/video with virtual backgrounds**, join from any modern browser, no install. See [Attendee Guide](https://www.plugnmeet.org/docs/user-guide/attendee).
- **Stable on weak networks:** Adaptive quality automatically lowers thumbnails and pauses off-screen videos to prevent lag.
- **HD screen sharing**, public + private chat with file sharing.

###### Collaborate live
- **Interactive whiteboard:** Draw together, upload PDFs/Office files and annotate them.
- **Shared notepad** for live meeting notes.
- **Embed any web content** (Google Doc, dashboard) or play a YouTube/file video in sync for all.

###### Engage and moderate
- **Breakout rooms** for small-group work, **polls with quiz mode**, raise hand, waiting room.
- **Granular permissions:** Lock mic, camera, chat, whiteboard per participant. See [Moderator Guide](https://www.plugnmeet.org/docs/user-guide/moderator).

###### Record and broadcast
- **Cloud recording to MP4** on your server + browser-based local recording.
- **Live stream to YouTube/Facebook** via RTMP, bring in studio feeds from [OBS via RTMP/WHIP](https://www.plugnmeet.org/blog/obs-rtmp-whip-ingress).
- **Telephone dial-in (SIP)** for audio-only join without internet.

###### Turn meetings into knowledge with AI
- Live captions, spoken translation, transcriptions, summaries, and chat assistant — powered by your own AI provider (OpenAI-compatible).
- You control when AI runs and where data goes. E2EE rooms automatically disable server-side AI that would need to hear audio. Details in [Security Overview](https://www.plugnmeet.org/docs/security-overview).

###### Measure with analytics
- **Deep meeting analytics & reports:** Attendance, talk time, engagement, polls, whiteboard and connection quality per session. Start with the [Meeting Analytics & Reports guide](https://www.plugnmeet.org/docs/others/analytics).
- Excel export via the PHP SDK, Moodle, WordPress and Joomla plugins.

###### Security and control
- Self-hosted data, HMAC-signed API + JWT sessions, TLS everywhere, optional **end-to-end encryption (E2EE)** for media and chat.
- Per-room analytics metadata (counts, not content) and full audit via API. Start with [Security & Privacy](https://www.plugnmeet.org/docs/security-overview).

##### How it works: 3 pieces

Forget complex diagrams. There are only three pieces:

1. **plugNmeet server — the engine.** Runs your conferences, manages rooms, users, recordings, and AI hooks. You either self-host it or use Cloud.
2. **Meeting interface — the room.** The modern web client your users see. Fully white-labeled via config: logo, colors, features on/off — no code fork needed.
3. **Connection — the bridge.** Links your website to the server:
   - **No-code:** WordPress / Joomla / Moodle plugin, or any LMS via LTI and BBB-compatibility.
   - **Pro-code:** PlugNmeet API + SDKs create rooms and mint join tokens from your backend. For iOS, Android, Flutter, React Native, or desktop, use the [hybrid mobile integration](https://www.plugnmeet.org/docs/developer-guide/mobile-app-integration) with LiveKit native SDKs + WebView UI.

```
Users (browser, 1 click, no install)
  |
Your site/app (WordPress, Moodle, Joomla, or custom code)
  |  creates room + join token via API
plugNmeet server (your own server or Cloud)
  |  hosts encrypted media + chat + recordings
Meeting room (your brand, your domain)
```

Your backend never exposes the secret: it signs the API request, gets a short-lived token for that user, and redirects them to the room.

##### Under the hood (for technical readers)

If you evaluate architecture, this is what matters:

- **Single Go binary** for the app logic — no runtime juggling, easy to deploy and scale horizontally.
- **LiveKit (SFU)** for media routing with simulcast/dynacast — one upload per sender, server forwards efficiently.
- **NATS JetStream** for signaling, chat, and collaboration state — with per-user, per-room least-privilege permissions and TLS.
- **Separate recorder** that outputs portable MP4. Run it on another machine so transcoding never slows live meetings.
- **Redis/NATS-KV (ephemeral) + MariaDB (history)** — live state disappears when the session ends; only what you opt to keep (recordings, analytics JSON, AI artifacts) is persisted.
- **API-first + BBB-compatible:** `POST /auth/*` with `API-KEY` + `HMAC-SHA256`, SDKs for [PHP](https://github.com/mynaparrot/plugNmeet-sdk-php) and [JavaScript](https://github.com/mynaparrot/plugNmeet-sdk-js), [webhooks](https://www.plugnmeet.org/docs/others/webhooks) and [hooks](https://www.plugnmeet.org/docs/others/hooks) for workflows, `/bigbluebutton/` compatibility for drop-in migration, and [hybrid mobile apps](https://www.plugnmeet.org/docs/developer-guide/mobile-app-integration) via LiveKit native SDKs + WebView (no proprietary mobile SDK required).

Deeper reads: [Benchmarks: how many users can a room hold](https://www.plugnmeet.org/docs/benchmark), [Backend Architecture Deep Dive](https://www.plugnmeet.org/blog/backend-architecture-deep-dive), [Scalable Deployment](https://www.plugnmeet.org/docs/developer-guide/scalable-setup), [Why we built plugNmeet](https://www.plugnmeet.org/blog/why-we-built-plugnmeet).

##### Why teams choose plugNmeet over SaaS or heavy stacks

- **Own your data and brand:** Your domain, your logo, your storage. No user lock-in, no per-user SaaS tax for basic meetings.
- **White-label without forking:** Change look and features via API/config and custom CSS. Updates do not break your theme.
- **Lightweight to run:** One binary + LiveKit + NATS replaces sprawling multi-language stacks. A small server handles testing; vertical scaling + isolated recorder covers hundreds of concurrent users before you ever need a cluster. See measured numbers in [Benchmarks](https://www.plugnmeet.org/docs/benchmark).
- **Portable recordings:** Plain MP4 files you can download, merge, or move — no proprietary post-processing maze.
- **Built for builders:** From one-click CMS plugins to headless `getClientFiles` embedding, [custom mobile apps](https://www.plugnmeet.org/docs/developer-guide/mobile-app-integration), and full automation, you are not forced into someone else's UX.

##### Get started in 15 minutes

###### Step 1: Get a server

**Option A — Self-host (free, full control)**

You need a clean Ubuntu/Debian server with public IP, 2 subdomains, and an email for SSL. Then:

```bash
wget https://raw.githubusercontent.com/mynaparrot/plugNmeet-install/main/install.sh
sudo bash install.sh
```

Save the Server URL, API Key, and Secret at the end. Full checklist in [Installation Quick Start](https://www.plugnmeet.org/docs/installation).

**Option B — Cloud (managed, fastest)**

No server work. Create an account on **[plugNmeet Cloud](https://www.plugnmeet.cloud)**, get credentials, and start with the [free plan](https://www.plugnmeet.cloud/pricing). Same privacy model, plus geo-routing for global audiences.

###### Step 2: Connect it to your site

- **WordPress / Joomla / Moodle?** Install the official plugin, paste server details, create your first room. See [WordPress](https://www.plugnmeet.org/docs/user-guide/wordPress-integration), [Joomla](https://www.plugnmeet.org/docs/user-guide/joomla-integration), [Moodle](https://www.plugnmeet.org/docs/user-guide/moodle-integration), [LTI](https://www.plugnmeet.org/docs/user-guide/lti).
- **Custom app?** Follow [PHP Quick Start](https://www.plugnmeet.org/docs/tutorials/quick_php) or read [API Intro](https://www.plugnmeet.org/docs/api/intro): create room → generate token → redirect user. No frontend SDK required to start.

##### Try it right now — live demo

Explore the UI and test the API without installing anything:

**Room demo:** [https://demo.plugnmeet.com/landing.html](https://demo.plugnmeet.com/landing.html)

**Demo API for developers:**
```
plugNmeet server URL: https://demo.plugnmeet.com
plugNmeet API KEY: plugnmeet
plugNmeet API SECRET: EXAMPLE_API_SECRET
```

##### Next steps

- **Just want to use it?** → [User Guide Overview](https://www.plugnmeet.org/docs/user-guide/overview) → [Moderator Guide](https://www.plugnmeet.org/docs/user-guide/moderator)
- **Want to host it?** → [Installation](https://www.plugnmeet.org/docs/installation) → [Security Overview](https://www.plugnmeet.org/docs/security-overview)
- **Want to build on it?** → [Developer Guide](https://www.plugnmeet.org/docs/developer-guide/intro) → [API Intro](https://www.plugnmeet.org/docs/api/intro) → [PHP Quick Start](https://www.plugnmeet.org/docs/tutorials/quick_php)
- **Need scale or migration?** → [Scalable Deployment](https://www.plugnmeet.org/docs/developer-guide/scalable-setup) → [Migrating from BBB](https://www.plugnmeet.org/docs/tutorials/migration-from-bbb)
- **Like the project?** Star us on [GitHub](https://github.com/mynaparrot/plugNmeet-server) and join [Discord](https://discord.gg/2X2ZaCHu4C).

<a id="page-35"></a>

## گزارش حضور، مشارکت و زمان گفتار

analytics هنگام ساخت اتاق با enable_analytics فعال می‌شود. فایل JSON پس از پایان جلسه به صورت MEETING_ANALYTICS در artifact ذخیره می‌شود و با رویداد analytics_proceeded قابل پیگیری است. room و users شامل رویدادهای نام‌دار هستند. زمان room بر حسب ثانیهٔ Unix و time در رویداد بر حسب میلی‌ثانیه است. talked و talked_duration شمارنده‌اند و به تنهایی بازهٔ شروع و پایان هر نوبت گفتار را فراهم نمی‌کنند. در آی‌نشست ثبت بازهٔ گفتار سمت سرور و تطبیق آن با زمان شروع هر ضبط لازم است. SDK PHP گزارش خام را به حضور و مجموع‌های قابل نمایش تبدیل می‌کند.

این بخش صفحهٔ `others/analytics` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/others/analytics) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/others/analytics.md) · [نسخهٔ محلی کامل](docs/upstream/docs/others/analytics.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Meeting Analytics & Deep Session Reports

plugNmeet's **meeting analytics** turns every session into a **detailed analytics report**: who joined and for how long, who spoke, who chatted, voted, raised a hand, shared a screen, or struggled with connection quality. Teams use it for **meeting reports** and attendance tracking; educators use it as a **classroom performance report** to measure participation equity, engagement trends, and technology barriers.

If `enable_analytics` was set in the room's `room_features` ([room create API](https://www.plugnmeet.org/docs/api/room/create)), plugNmeet aggregates analytics events in memory during the session and writes a single JSON file when the session ends. The file is stored as an artifact of type `MEETING_ANALYTICS`.

For the why behind these metrics — engagement questions, equitable participation, and privacy-by-default design — see [From Black Box to Blueprint: Using Analytics to Understand Your Online Classroom](https://www.plugnmeet.org/blog/classroom-analytics-for-educators).

##### What you can build from this data

- **Meeting attendance report:** per-user `duration` from `joined` / `left` timelines, plus `room_total_users` and `room_duration`.
- **Participant engagement analytics:** `talked` / `talked_duration`, `mic_status`, `webcam_status`, `screen_share_status`, `raise_hand`, `voted_poll`, `public_chat` / `private_chat`.
- **Classroom performance report:** combine talk time, poll votes, whiteboard annotation (`whiteboard_annotated`), chat files, and hand raises to compare participation across students or sessions. The Moodle plugin does this out of the box with cumulative attendance and activity completion rules — see [Moodle integration](https://www.plugnmeet.org/docs/user-guide/moodle-integration).
- **Connection quality report:** `connection_quality` (`excellent` / `good` / `poor`) correlated with early `left` events to distinguish disengagement from network issues.
- **Excel / LMS exports:** the Moodle, WordPress, and Joomla plugins format the same data into per-room and per-user sheets. See [WordPress artifacts](https://www.plugnmeet.org/docs/user-guide/wordPress-integration) for an example.

Analytics is **opt-in and privacy-preserving**: disabled by default, in-memory only, and permanently discarded if `enable_analytics` is false. Only metadata and counters are stored — never chat text, audio, or video content. See [Security overview](https://www.plugnmeet.org/docs/security-overview).

##### How to get the analytics file

1. Enable it at room creation: `"enable_analytics": true` in `room_features`.
2. Wait for the session to end. The server exports the JSON and emits the webhook event `analytics_proceeded` with the `artifact_id` — see [Webhooks](https://www.plugnmeet.org/docs/others/webhooks).
3. List files with [`/artifact/fetch`](https://www.plugnmeet.org/docs/api/artifact/fetch) using `"type": "MEETING_ANALYTICS"`.
4. Get a one-time token from [`/artifact/getDownloadToken`](https://www.plugnmeet.org/docs/api/artifact/download), then download:
   `GET https://Your-Server/download/artifact/<TOKEN>`

The deprecated `/analytics/fetch` and `/analytics/getDownloadToken` endpoints still work but will be removed. Use `/artifact/*` for new code.

##### Raw analytics JSON format

The raw file is intentionally low level — a complete event log, not a ready-made report. Don't aggregate it by hand; use the [PHP SDK formatter](https://www.plugnmeet.org/docs/others/analytics#format-with-the-php-sdk) below.

Top level is always:

```json
{
  "room": {},
  "users": []
}
```

###### Room object

Base fields are written directly (see `AnalyticsRoomInfo` and `exportAnalyticsToJSON` in the server):

| Field | Type | Notes |
| ----- | ---- | ----- |
| `room_id` | string | Your `room_id`. |
| `room_title` | string | Room title at end time. |
| `room_total_users` | int | Number of distinct users seen. |
| `room_creation` | int | Unix seconds. |
| `room_ended` | int | Unix seconds. |
| `room_duration` | int | Seconds, `room_ended - room_creation`. |
| `enabled_e2ee` | bool | From room metadata. |
| `events` | array | Low-level event list, see below. |

###### User object

```json
{
  "user_id": "111",
  "name": "Test User",
  "is_admin": false,
  "ex_user_id": "111",
  "events": []
}
```

`user_id` is the internal LiveKit participant ID, `ex_user_id` is the `userId` you passed in the join token (empty if `auto_gen_user_id` was used).

###### Event object

Both room and user events share the same shape:

```json
{
  "name": "joined",
  "total": 6,
  "values": [{ "time": 1500000, "value": "1500000" }]
}
```

- `name` is the proto `AnalyticsEvents` name lowercased with the `ANALYTICS_EVENT_ROOM_` / `ANALYTICS_EVENT_USER_` prefix stripped, e.g. `poll_added`, `recording_status`, `joined`, `mic_status`. This stripping is done in `processEventKey`.
- `time` is Unix **milliseconds**.
- `total` / `values` encoding depends on how the server stored the event in Redis:
  - `HSET` (timestamped events like `joined`, `mic_status`): `total` = number of entries, `values` = full list.
  - `INCRBY` (counters like `talked`, `public_chat`): `total` = counter, `values` is empty.
  - `SET` (string values): `total` = `1`, `values` has one entry.

Known event names (from `plugnmeet_analytics.proto`):

Room: `poll_added`, `recording_status`, `rtmp_status`, `breakout_room`, `shared_notepad_status`, `external_display_link_status`, `external_media_player_status`, `speech_service_status`, `speech_service_total_usage`, `ingress_created`, `whiteboard_files`, insights/transcription/translation counters.

User: `joined`, `left`, `mic_status`, `talked`, `talked_duration` (ms), `webcam_status`, `screen_share_status`, `voted_poll`, `public_chat`, `private_chat`, `chat_files`, `whiteboard_annotated`, `whiteboard_files`, `speech_services_usage`, `speech_services_status`, `interface_visibility`, `raise_hand`, `connection_quality`, `reaction`.

Status values inside `mic_status` / `webcam_status` / room status events use prefixes like `ANALYTICS_STATUS_STARTED`, `ANALYTICS_STATUS_ENDED`, `ANALYTICS_STATUS_MUTED`, or `excellent` / `good` / `poor` for connection quality.

**Fields evolve with new features**
New plugNmeet features can add or extend room fields, user fields, and event names at any time, so the lists above are a snapshot, not a fixed schema. Before parsing strictly or validating against hardcoded field names, always check the release notes and the authoritative [`plugnmeet_analytics.proto`](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_analytics.proto) file. When using the PHP SDK, prefer `getRoomFields()` / `getUserFields()` over hardcoded lists so your code picks up new fields automatically.


##### Format with the PHP SDK

The [`plugNmeet-sdk-php`](https://github.com/mynaparrot/plugNmeet-sdk-php) package ships `Mynaparrot\Plugnmeet\AnalyticsFormatter`, which converts the raw JSON into an attendance-ready meeting report. The Moodle, WordPress, and Joomla plugins use the same class to generate their Excel exports.

```bash
composer require mynaparrot/plugnmeet-sdk
```

###### 1. Fetch and download the JSON

```php
<?php
require __DIR__ . '/vendor/autoload.php';

use Mynaparrot\Plugnmeet\PlugNmeet;
use Mynaparrot\PlugnmeetProto\FetchArtifactsReq;
use Mynaparrot\PlugnmeetProto\GetArtifactDownloadTokenReq;
use Mynaparrot\PlugnmeetProto\RoomArtifactType;

$pnm = new PlugNmeet('https://your-server', 'api-key', 'api-secret');

// list analytics artifacts for a room
$req = new FetchArtifactsReq();
$req->setRoomIds(['room01']);
$req->setType(RoomArtifactType::MEETING_ANALYTICS);
$res = $pnm->fetchArtifacts($req);

$artifactId = $res->getResult()->getArtifactsList()[0]->getArtifactId();

// download
$tokenReq = new GetArtifactDownloadTokenReq();
$tokenReq->setArtifactId($artifactId);
$token = $pnm->getArtifactDownloadToken($tokenReq)->getToken();

$json = file_get_contents('https://your-server/download/artifact/' . $token);
$raw = json_decode($json, true);
```

###### 2. Format it into a report

```php
use Mynaparrot\Plugnmeet\AnalyticsFormatter;

$formatter = new AnalyticsFormatter($raw, 'Asia/Dhaka');
$data = $formatter->getFormattedEventData();

$room  = $data['room'];   // associative array
$users = $data['users'];  // list of associative arrays

// helpers
echo $formatter->formatTimestamp($raw['room']['room_creation'], false);
echo $formatter->formatSecondsToTime($room['room_duration']);
print_r($formatter->getRoomFields());
print_r($formatter->getUserFields());
```

The wrapper in the SDK examples does the same in one call:

```php
$formatter = plugNmeetConnect::getAnalyticsFormatter($json, $userTimezone);
```

Invalid timezones fall back to `UTC`.

###### What the formatter produces

Room fields (`getRoomFields()`):

```
room_id, room_title, room_creation, room_ended, room_duration,
room_total_users, enabled_e2ee, recording_status, rtmp_status,
shared_notepad_status, ingress_created, breakout_room,
external_display_link_status, external_media_player_status
```

Plus `polls` and `whiteboard_files` when present. `room_creation` / `room_ended` become `d-m-Y H:i:s P` strings in the given timezone. Status events become start-counts.

User fields (`getUserFields()`):

```
name, ex_user_id, is_admin, duration, joined, left, mic_status,
mic_muted, mic_duration, talked, talked_duration, webcam_status,
webcam_duration, raise_hand, voted_poll, whiteboard_annotated,
whiteboard_files, screen_share_status, public_chat, private_chat,
chat_files, interface_invisible, connection_quality
```

Key transformations:

- `joined` / `left`: sorted lists of Unix-ms timestamps.
- `duration`: total seconds present, computed from a join/left timeline. Duplicate joins are ignored, unmatched leaves are ignored, an unfinished session extends to `room_ended`, events outside `[room_creation, room_ended]` are ignored.
- `mic_duration` / `webcam_duration`: seconds between `ANALYTICS_STATUS_STARTED` and `ANALYTICS_STATUS_ENDED` (prefix match, so `ANALYTICS_STATUS_STARTED:node_01` counts). An open device session closes at user `left` or room end. Duplicate starts are ignored.
- `talked_duration`: raw `total` is milliseconds → converted to seconds with `ceil(ms / 1000)`.
- `talked`, `public_chat`, `private_chat`, `chat_files`, `whiteboard_annotated`, `raise_hand`: plain counters from `total`.
- `interface_visibility` → `interface_invisible` (count of `hidden`).
- `connection_quality` → `{ excellent, good, poor }` counts.
- `voted_poll`: per-user vote count; also increments `room.polls[poll_id].options[option_id].responses`.
- `duration`, `talked_duration`, `mic_duration`, `webcam_duration` are clamped to `min(room_duration, room_ended - room_creation)`. If room metadata is missing, the raw value is kept.

Formatted example — one row of a meeting analytics report:

```json
{
  "room": {
    "room_id": "room01",
    "room_title": "Test Room",
    "room_creation": "17-05-2026 10:00:00 +06:00",
    "room_duration": 3000,
    "recording_status": 1
  },
  "users": [
    {
      "name": "Test User",
      "duration": 1200,
      "mic_duration": 500,
      "talked_duration": 82,
      "connection_quality": { "excellent": 10, "good": 2, "poor": 0 }
    }
  ]
}
```

Use `formatSecondsToTime()` to render duration fields as `HH:MM:SS` for tables or CSV/XLSX export.

##### Related resources

- [Classroom analytics for educators](https://www.plugnmeet.org/blog/classroom-analytics-for-educators) — pedagogical guide to engagement analytics.
- [Moodle integration](https://www.plugnmeet.org/docs/user-guide/moodle-integration) — cumulative classroom performance reports with attendance and activity completion.
- [Fetch artifacts API](https://www.plugnmeet.org/docs/api/artifact/fetch) and [download token API](https://www.plugnmeet.org/docs/api/artifact/download).
- [Webhooks](https://www.plugnmeet.org/docs/others/webhooks) — `analytics_proceeded` event.
- [Security overview](https://www.plugnmeet.org/docs/security-overview) — opt-in analytics and data retention.

<a id="page-36"></a>

## Hookهای اسکریپت و ذخیره‌سازی

hookها در سرور و recorder برای انتقال فایل و پردازش سفارشی اجرا می‌شوند. فرایند طولانی‌عمر یا فرمان یک‌باره باید برای هر ورودی یک خط JSON روی stdout برگرداند و لاگ را به stderr بفرستد. اجرای hook مسدودکننده است و پاسخ‌ندادن، مسیر سرویس را متوقف می‌کند. چند hook زنجیره‌ای ورودی و خروجی را منتقل می‌کنند. فعال‌کردن hook فایل، مسئولیت پاک‌سازی و سازگاری مسیر را به پیاده‌سازی منتقل می‌کند. حذف محلی فقط پس از انتقال موفق و طبق نقش hook انجام شود.

این بخش صفحهٔ `others/hooks` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/others/hooks) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/others/hooks.md) · [نسخهٔ محلی کامل](docs/upstream/docs/others/hooks.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Scripting & Storage Hooks

plugNmeet features a powerful hooking system that allows you to execute custom scripts or commands at key points in the media and file management lifecycle. This enables deep customization, letting you integrate with any external storage provider (e.g., S3, Google Cloud), call custom APIs, or orchestrate complex, multi-server pipelines.

Hooks are available in both the `server` and `recorder` components.

**Looking for Examples?**
This page covers the technical details of the hooking system. For practical, real-world examples in various programming languages, check out our blog posts tagged with **[hooks](https://www.plugnmeet.org/blog/tags/hooks)**.


**Source of Truth**
The JSON data structures for hooks may change over time. This documentation provides examples, but for the most up-to-date definitions, please refer to the official `plugnmeet-protocol` repository:

[https://github.com/mynaparrot/plugnmeet-protocol/tree/main/hooks](https://github.com/mynaparrot/plugnmeet-protocol/tree/main/hooks)


##### Core Concepts

All hooks, regardless of where they run, share the same fundamental design. Understanding these core concepts is essential before implementing your own custom scripts.

###### Execution Models

For maximum performance and flexibility, hook scripts can be configured as either **long-lived processes** or **one-shot commands**.

*   **Long-Lived Processes**: When a plugNmeet component (like `server` or `recorder`) starts, it launches your script **once**. The script then runs continuously, listening for requests. This model is highly efficient as it avoids the overhead of starting a new process for every event. **This is the recommended approach for most use cases.**
*   **One-shot Commands**: These are simple commands (like `curl`, `wget`, or the built-in `http-request`) that are executed for each hook event. They are suitable for simple, self-contained tasks.

###### Communication Protocol

Communication with your hook script happens over standard I/O pipes using **newline-delimited JSON**.

*   **`stdin`**: Your script reads requests from `stdin`. Each line is a complete JSON object representing a single request.
*   **`stdout`**: For each request received, your script **must** print a single line of JSON to `stdout` as its response.
*   **`stderr`**: You can use `stderr` for logging and debugging within your script. This output is ignored by plugNmeet but is invaluable for development.

**CRITICAL: Scripts are Blocking**
The call to a hook script is **synchronous and blocking**. Your script **MUST** write a response to `stdout` for every request it receives. If a script fails to return a response, the plugNmeet service will hang indefinitely.


###### The Pipeline Model

If you define multiple scripts for a single hook, they form a pipeline. The `stdout` response from the first script becomes the `stdin` request for the second, and so on.

If a script in the chain doesn't need to modify the data (e.g., it only logs the event), it **must** still pass the original, unmodified JSON object through to `stdout`.

###### Error & Data Handling

*   **Valid JSON Response**: If your script returns valid JSON, it is passed to the next script or used as the final result.
*   **Invalid JSON Response**: If the response is not valid JSON, plugNmeet logs a **warning** and discards the output. The *original* JSON data (the `stdin` to your script) is passed to the next script in the pipeline. This prevents a single faulty script from breaking the entire chain.
*   **Reporting Errors**: If your script encounters an error, it should populate the `error` field in its JSON response. It is crucial to **always return the complete input JSON object, with the `error` field populated**, to ensure subsequent scripts in a pipeline receive the expected data structure. The main application will log this error.

###### How to Create a Hook

1.  **Create a Script or Command**: Write your logic in any language (Shell, NodeJS, Go, etc.).
2.  **Make it Executable**: Ensure your script has execute permissions (e.g., `chmod +x your_script.js`).
3.  **Configure in `config.yaml`**: Add the script's absolute path or the command to the appropriate `hooks` section in your `config.yaml`.

###### Built-in `http-request` Utility

plugNmeet provides a convenient one-shot command, `http-request`, for sending the hook's JSON payload to an HTTP/HTTPS endpoint.

**Usage:**
`http-request <URL>`

**Example in `config.yaml`:**
```yaml
scripts:
  - script: "http-request http://localhost:8090/your/endpoint"
    is_one_shot: true
```

---

##### Key Responsibilities

###### File Cleanup Responsibility

**Important**
When you enable a hook that receives a file path (like `upload_hook` or `post_transcoding`), plugNmeet **disables its own automatic file cleanup** for that stage.


This is a critical safety feature. The application delegates file management to your script because it cannot know its nature (e.g., whether it's an uploader or just a notifier). If the application deleted a file before your script could process it, it would cause an error.

Therefore, the responsibility for file management is transferred to you, and your script's role determines its responsibility:

*   **If your script MOVES or UPLOADS the file** (e.g., to S3), it **MUST** delete the local source file from `input_path` after the transfer is successful. This is essential to prevent your server's disk from filling up.

*   **If your script only OBSERVES the file** (e.g., for logging, analytics, or sending a notification) and does not move it, it **MUST NOT** delete the file. The file is still needed by plugNmeet for its own internal storage or for subsequent hooks in a pipeline.

###### Path Consistency Responsibility

**Path Consistency is Your Responsibility**
**plugNmeet does not validate the `output_path` you return from a hook. It is stored as a string and used as the `input_path` for subsequent `download_hook` and `delete_hook` calls.**


*   **If your script modifies `output_path`** (e.g., changing a local path to an S3 key in an `upload_hook` or `post_transcoding` hook), you take full responsibility for that path. You **MUST** also implement a corresponding `download_hook` and `delete_hook` that can understand and process the custom path format you have defined.

*   **If your script is only for observation** (e.g., logging stats or sending a notification) and does not modify the `output_path`, then you do not need to provide the other hooks. The default workflow will continue with the original path.

Failing to provide compatible `download_hook` and `delete_hook` scripts after changing the `output_path` will result in broken downloads and deletions.


---

##### Server Storage Hooks (`server`)

Server hooks allow you to override the default local file storage for room artifacts, chat files, and recordings, enabling integration with any external storage provider.

**Configuration in `server/config.yaml`:**
```yaml
hooks:
  # 'pool_size' controls parallel execution. Default: 1.
  # 'hook_timeout' sets a timeout. Default: 5m.
  upload_hook:
    pool_size: 2
    scripts:
      - script: "/path/to/your/upload_script.sh"
        is_one_shot: false

  download_hook:
    scripts:
      - script: "/path/to/your/download_script.sh"
        is_one_shot: false
  # ... other hooks ...
```

###### Hook Types

###### `upload_hook`
*   **Purpose**: Upload a file or directory (e.g., room analytics, whiteboard images) to external storage.
*   **Input (`UploadHookData`)**:
    ```json
    {
      "input_path": "/path/on/server/disk/analytics.json",
      "input_directory_path": "", // or "/path/to/converted/images/"
      "hook_file_type": "artifact",
      "room_id": "room01",
      "room_sid": "SID_d82k3s9d2l"
    }
    ```
*   **Script's Job**: Upload the content from `input_path` or `input_directory_path`. Return the JSON with `output_path` set to a unique storage identifier. **After a successful upload, you must delete the local source file/directory.**

###### `download_hook`
*   **Purpose**: Provide a secure way for users to download a file from external storage.
*   **Input (`DownloadHookData`)**:
    ```json
    {
      "input_path": "artifacts/room01/analytics.json", // The storage identifier
      "hook_file_type": "artifact"
    }
    ```
*   **Script's Job**: Return JSON specifying an `action`.
    *   **`action: "redirect"` (Recommended)**: Generate a temporary, pre-signed URL and return it in the `redirect_url` field.
    *   **`action: "serve_local"`**: Download the file to a temporary local path on the server and return the path in `output_path` and the file's `mime_type`.

###### `delete_hook`
*   **Purpose**: Delete a file from external storage.
*   **Input (`DeleteHookData`)**:
    ```json
    {
      "input_path": "artifacts/room01/analytics.json", // The storage identifier
      "hook_file_type": "artifact"
    }
    ```
*   **Script's Job**: Delete the file from storage and return a confirmation response.

###### `resumable_upload_hook`
*   **Purpose**: Handle chunked file uploads for the chat feature, typically by offloading logic to a service like S3's Multipart Upload.
*   **Input (`ResumableUploadHookData`)**: Contains a `type` field (`part-check`, `part-upload`, `merge`) that dictates the required action.
*   **Script's Job**: Interact with your storage backend to check, upload, or merge file chunks. The response must indicate the outcome (e.g., `part_exists`, `merge_success`).

###### `room_end_hook`
*   **Purpose**: Perform cleanup tasks after a room session has ended (e.g., clean up abandoned resumable upload chunks).
*   **Input (`RoomEndHookData`)**:
    ```json
    {
      "room_id": "room01",
      "room_sid": "SID_d82k3s9d2l"
    }
    ```
*   **Script's Job**: Perform cleanup and return a confirmation message.

---

##### Recorder Hooks (`recorder`)

Recorder hooks are used to manage the recording file as it moves through the transcoding pipeline. This is essential for multi-server deployments where recording and transcoding may happen on different machines.

**Configuration in `recorder/config.yaml`:**
```yaml
hooks:
  # 'pool_size' controls how many hook pipelines can run in parallel. Default: 1.
  # 'hook_timeout' sets a timeout for the entire hook chain. Default: 1h.
  post_recording:
    pool_size: 2
    hook_timeout: 2h
    scripts:
      - script: "./scripts/post-recording/upload.sh"
        is_one_shot: false

  pre_transcoding:
    scripts:
      - script: "./scripts/pre-transcoding/download.sh"
        is_one_shot: false

  post_transcoding:
    scripts:
      - script: "./scripts/post-transcoding/notify.sh"
        is_one_shot: false
```

###### Data Payload: `RecordingHookData`

All recorder hooks receive and are expected to return a JSON object with this structure.

```json
{
  "task": "single",
  "recording_id": "REC_ax9s3djn2s",
  "room_table_id": 123,
  "room_id": "room01",
  "room_sid": "SID_d82k3s9d2l",
  "file_name": "REC_ax9s3djn2s.mp4",
  "recorder_id": "node_01",
  "file_size": 123.45,
  "input_path": "/path/to/file.mp4",
  "input_paths": [],
  "output_path": "",
  "error": "",
  "should_cleanup": false,
  "source_for_cleanup": ""
}
```

###### Hook Stages

###### 1. `post_recording`
*   **When**: Runs on the **RECORDER** node after a raw recording file is saved.
*   **Purpose**: Upload the raw file to a network-accessible location (e.g., S3, NFS) so it can be accessed by a transcoder.
*   **Input**: The `input_path` field contains the path to the raw recording file on the recorder's local disk.
*   **Script's Job**: Upload the file and return the JSON with `output_path` set to the file's new network location or identifier (e.g., an S3 key). **After a successful upload, you must delete the local source file from `input_path`.**

###### 2. `pre_transcoding`
*   **When**: Runs on the **TRANSCODER** node before `ffmpeg` processing begins.
*   **Purpose**: Download the raw file from a network location to a temporary local path on the transcoder.
*   **Input**: The `input_path` field contains the network location/identifier of the raw file (this is the `output_path` from the `post_recording` stage).
*   **Script's Job**: Download the file and return the JSON with `output_path` set to the new **local path** on the transcoder's disk.

###### 3. `post_transcoding`
*   **When**: Runs on the **TRANSCODER** node after `ffmpeg` successfully creates the final `.mp4` file.
*   **Purpose**: Upload the final processed file to permanent storage and perform cleanup.
*   **Input**: The `input_path` field contains the path to the final, transcoded `.mp4` file on the local disk of the transcoder node.
*   **Script's Job**: Upload the final file and return the JSON, optionally updating `output_path`. The `should_cleanup` and `source_for_cleanup` fields can be used to manage cleanup of temporary files from the `pre_transcoding` stage.

---

##### Example: Long-Lived Node.js Script

This example shows the basic structure for a long-lived script that safely parses JSON, performs an action, and returns a response.

```javascript
#!/usr/bin/env node
// scripts/my_hook.js

const readline = require('readline');
const fs = require('fs');

// Use stderr for logging so it doesn't interfere with stdout
const log = (message) => {
  console.error(`[MyHook] ${new Date().toISOString()}: ${message}`);
};

log('Starting long-lived hook script...');

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
  terminal: false,
});

rl.on('line', (line) => {
  let requestData;
  try {
    requestData = JSON.parse(line);
    log(`Received request for room: ${requestData.room_id || 'N/A'}`);

    // --- Your Custom Logic Here ---
    // e.g., upload to S3, call an API, etc.
    // After a successful upload, remember to clean up the source file.
    if (requestData.input_path) {
      // In a real script, you would do this AFTER a successful upload.
      // fs.unlinkSync(requestData.input_path);
      // log(`Cleaned up ${requestData.input_path}`);
    }

    requestData.processed_by_hook = true;
    // ---

    // ALWAYS write a valid JSON response to stdout
    process.stdout.write(JSON.stringify(requestData) + '\n');

  } catch (e) {
    log(`ERROR: ${e.message}`);
    // If an error occurs, return the original requestData object with an 'error' field.
    // This ensures the full structure is maintained for subsequent scripts.
    const errorResponse = requestData 
      ? { ...requestData, error: e.message, output_path: "" } // Clear output_path on error
      : { error: `Failed to parse incoming JSON: ${e.message}` };
      
    process.stdout.write(JSON.stringify(errorResponse) + '\n');
  }
});

rl.on('close', () => {
  log('Stdin closed. Exiting script.');
  process.exit(0);
});
```

<a id="page-37"></a>

## وب‌هوک‌ها و چرخهٔ جلسه

وب‌هوک POST با JSON ارسال می‌شود و Authorization و Hash-Token شامل JWT امضاشده با هش SHA256 بدنه هستند. room.room_id برای ارتباط با جلسه و شناسهٔ نشست برای تفکیک اجراها ثبت می‌شود. session_ended پایان نشست را پوشش می‌دهد، حتی اگر کاربری وارد نشده باشد. end_recording توقف ضبط است و recording_proceeded آماده‌شدن فایل را اعلام می‌کند. رویداد ممکن است پس از پاک‌سازی جلسه برسد. گیرندهٔ آی‌نشست باید رویداد را پایدار و idempotent پردازش کند و فایل جامانده را از API بازیابی کند.

این بخش صفحهٔ `others/webhooks` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/others/webhooks) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/others/webhooks.md) · [نسخهٔ محلی کامل](docs/upstream/docs/others/webhooks.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Webhooks

Plug-N-Meet can notify your application about various events by sending webhook requests to a specified URL. You can configure this URL in the server settings or provide it when creating a room.

##### Receiving Webhooks

Webhook requests are sent as HTTP POST requests to the URLs you have configured, either in the server's `config.yml` file or during room creation. Each webhook event is encoded as JSON and included in the request body.

The request will have the `Content-Type` header set to `application/json`. Make sure your web server is configured to accept payloads with this content type.

Plug-N-Meet uses the same security pattern as LiveKit. To verify that webhook requests originate from Plug-N-Meet, each request includes `Authorization` and `Hash-Token` headers containing a signed JWT token. The token includes a SHA256 hash of the payload. For a PHP example, see [webhook.php](https://github.com/mynaparrot/plugNmeet-sdk-php/blob/main/examples/webhook.php).

Every event contains `event`, `room`, `id`, and `createdAt`. `room` is a `NotifyEventRoom` with at least `room_id` and `sid`. Always match sessions by `room.room_id`. If LiveKit fails to include room info in its callback, plugNmeet drops it and never forwards it.

##### Session lifecycle: which event to trust

A plugNmeet session is not the same as a LiveKit media room. `room_created` and `session_ended` are sent by plugNmeet itself and always fire. `room_started` and `room_finished` are relayed from LiveKit and only fire if at least one participant joined and LiveKit notified us.

- Use `room_created` + `session_ended` for reliable session tracking, cleanup, and billing.
- Use `room_started` + `room_finished` only for media state (e.g. someone actually joined, media room closed).

| Event | Sent by | When it fires |
|---|---|---|
| `room_created` | plugNmeet | Right after `CreateRoom` API succeeds. Always fires. |
| `room_started` | LiveKit, relayed by plugNmeet | When the LiveKit room becomes active, i.e. first participant joined. Only fires if someone joined and LiveKit notified us. |
| `room_finished` | LiveKit, relayed by plugNmeet | When the LiveKit room closes. Only fires if the room had started and LiveKit notified us. Missing if nobody ever joined. |
| `session_ended` | plugNmeet | During `EndRoom` flow, whether triggered by API `EndRoom` call or by janitor (empty timeout, duration exceeded). Always fires, even if nobody ever joined. Sent before the LiveKit room is torn down. |

If you end a session via API or janitor ends an empty session, you will always get `session_ended`. You will only get `room_finished` afterwards if the LiveKit room had actually started.

##### Events

You can review the event definitions [here](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_common.proto#L10).

All webhook events include the following fields:
- `id`: A UUID identifying the event
- `createdAt`: UNIX timestamp in seconds
- `room`: Contains at least `room_id`, `sid`, `creation_time`, and `metadata`

###### Room Created

Sent by plugNmeet on session creation. Use it to create your local session record.

```js
interface CommonNotifyEvent {
  event: 'room_created'
  room: Room
}
```

###### Room Started

Relayed from LiveKit when the media room starts. Fires only if at least one user joined. plugNmeet enriches it with `sid` and `metadata` before forwarding.

```js
interface CommonNotifyEvent {
  event: 'room_started'
  room: Room
}
```

###### Room Finished

Relayed from LiveKit when the media room closes. Fires only if the room had started. If the room was ended via plugNmeet API, this may arrive after `session_ended`.

```js
interface CommonNotifyEvent {
  event: 'room_finished'
  room: Room
}
```

###### Session Ended

Special plugNmeet-only event. Sent from `onAfterRoomEnded` (`pkg/models/room_end.go`) before LiveKit teardown, whenever a session is wiped by API call or janitor. Always fires. Use it as the authoritative end-of-session signal for cleanup and billing. Payload is minimal: `room_id`, `sid`, `metadata`, `creation_time`.

```js
interface CommonNotifyEvent {
  event: 'session_ended'
  room: Room
}
```

###### Participant Joined

Relayed from LiveKit when a participant joins. plugNmeet overwrites `sid` and `metadata` from its own store before forwarding.

```js
interface CommonNotifyEvent {
  event: 'participant_joined'
  room: Room
  participant: ParticipantInfo
}
```

###### Participant Left

Relayed from LiveKit when a participant leaves. Still sent for internal users (SIP, ingress, TTS agents).

```js
interface CommonNotifyEvent {
  event: 'participant_left'
  room: Room
  participant: ParticipantInfo
}
```

###### Track Published

Relayed from LiveKit when a mic, camera, or screen-share track is published.

```js
interface CommonNotifyEvent {
  event: 'track_published'
  room: Room
  participant: ParticipantInfo
  track: TrackInfo
}
```

###### Track Unpublished

Relayed from LiveKit when a track is unpublished/muted by stopping publication.

```js
interface CommonNotifyEvent {
  event: 'track_unpublished'
  room: Room
  participant: ParticipantInfo
  track: TrackInfo
}
```

###### Recording Started

Sent by the recorder service when recording starts. Contains `recording_info` with record and recorder IDs.

```js
interface CommonNotifyEvent {
  event: 'start_recording'
  room: Room
  recording_info: RecordingInfoEvent
}
```

###### Recording Ended

Sent by the recorder service when recording stops. File may not be ready yet; wait for `recording_proceeded` for the downloadable file.

```js
interface CommonNotifyEvent {
  event: 'end_recording'
  room: Room
  recording_info: RecordingInfoEvent
}
```

###### Recording Proceeded

Sent when the recording file is processed and ready. This may arrive after session cleanup.

```js
interface CommonNotifyEvent {
  event: 'recording_proceeded'
  room: Room
  recording_info: RecordingInfoEvent
}
```

###### RTMP Started

Sent by the recorder service when RTMP livestreaming starts.

```js
interface CommonNotifyEvent {
  event: 'start_rtmp'
  room: Room
  recording_info: RecordingInfoEvent
}
```

###### RTMP Ended

Sent by the recorder service when RTMP livestreaming stops.

```js
interface CommonNotifyEvent {
  event: 'end_rtmp'
  room: Room
  recording_info: RecordingInfoEvent
}
```

###### Artifact Created

This single, unified event is triggered whenever a new artifact is generated by the system, such as a transcription or a meeting summary. It replaces multiple older, individual webhooks with a more flexible and scalable model.

For details on the payload structure and possible artifact types, please refer to the [protobuf definition](https://github.com/mynaparrot/plugnmeet-protocol/blob/main/proto_files/plugnmeet_room_artifacts.proto#L91).

```js
interface CommonNotifyEvent {
  event: 'artifact_created'
  room: Room
  room_artifact: RoomArtifactWebhookEvent
}
```

###### Analytics Proceeded

Sent when the analytics export for a session is ready. Contains `analytics.file_id` for download.

```js
interface CommonNotifyEvent {
  event: 'analytics_proceeded'
  room: Room
  analytics: AnalyticsEvent
}
```

---

##### Example Webhook Handler

The following pseudo-code demonstrates how to handle incoming webhooks. Your application should inspect the `event` property to determine what action to take. The `artifact_created` event is special, as it requires checking a nested `type` field to understand the payload.

```javascript
// A mapping from the enum number to a human-readable string.
// See the protobuf definition for the complete, up-to-date list.
const ArtifactType = {
  0: 'UNKNOWN_ARTIFACT',
  2: 'MEETING_SUMMARY',
  4: 'SPEECH_TRANSCRIPTION',
  5: 'SPEECH_TRANSCRIPTION_USAGE',
};

function handleWebhook(payload) {
  console.log(`Received event: ${payload.event}`);
  // Always match by room_id. Events without room.room_id are never sent.
  const roomId = payload.room?.room_id;

  switch (payload.event) {
    case 'room_created':
      // Reliable session start.
      console.log(`Session '${roomId}' created.`);
      break;

    case 'room_started':
      // Media room started (someone joined).
      console.log(`Room '${roomId}' has started.`);
      break;

    case 'participant_joined':
      console.log(`User '${payload.participant.name}' joined room '${roomId}'.`);
      break;

    case 'artifact_created':
      handleArtifactCreated(payload);
      break;

    case 'room_finished':
      console.log(`Room '${roomId}' has finished.`);
      break;

    case 'session_ended':
      // Authoritative session end. Always fires via API or janitor.
      console.log(`Session '${roomId}' has ended.`);
      break;

    default:
      console.log(`Received an unhandled event: ${payload.event}`);
  }
}

// A dedicated handler for the 'artifact_created' event
function handleArtifactCreated(payload) {
  const artifact = payload.room_artifact;
  const artifactType = ArtifactType[artifact.type] || 'UNKNOWN_ARTIFACT';
  
  console.log(`A new artifact of type '${artifactType}' was created.`);

  switch (artifactType) {
    case 'MEETING_SUMMARY':
    case 'SPEECH_TRANSCRIPTION':
      // These are file-based artifacts.
      if (artifact.metadata && artifact.metadata.file_info) {
        const fileInfo = artifact.metadata.file_info;
        console.log(`File available at: ${fileInfo.file_path}`);
        // Your logic here: e.g., download the file, link it to a recording.
      }
      break;
      
    case 'SPEECH_TRANSCRIPTION_USAGE':
      // This is a usage-based artifact.
      if (artifact.metadata && artifact.metadata.duration_usage) {
        const usage = artifact.metadata.duration_usage;
        console.log(`Transcription usage logged: ${usage.duration_sec} seconds.`);
        // Your logic here: e.g., update billing records.
      }
      break;
      
    default:
      console.log(`Received an unhandled artifact type: ${artifactType}`);
  }
}

// Example incoming payload for 'artifact_created' (a meeting summary)
const samplePayload = {
  "event": "artifact_created",
  "id": "EV_abc123",
  "createdAt": 1761815000,
  "room": {
    "sid": "RM_zG9e4bY2pD-1761814595173",
    "room_id": "room01"
  },
  "room_artifact": {
    "type": 2, // Corresponds to MEETING_SUMMARY
    "artifact_id": "ART_xyz789",
    "metadata": {
      "file_info": {
        "file_path": "/path/to/summary.txt",
        "file_size": 1024,
        "mime_type": "text/plain"
      }
    }
  }
};

// You would call your handler with the body of the POST request.
handleWebhook(samplePayload);
```

<a id="page-38"></a>

## امنیت، حریم خصوصی و چرخهٔ داده

درخواست مدیریتی با HMAC-SHA256 احراز می‌شود و کاربر با توکن کوتاه‌عمر وارد اتاق می‌شود. رسانهٔ WebRTC در انتقال رمزگذاری می‌شود. E2EE حالت جداگانه‌ای با پیامدهای مستقیم برای ضبط و سرویس‌های هوش مصنوعی دارد و نحوهٔ تحویل کلید به عامل مجاز باید بررسی شود. وضعیت بلادرنگ، دادهٔ موقت، اطلاعات پایدار دیتابیس، ضبط و artifact چرخهٔ نگهداری متفاوت دارند. دادهٔ ارسالی به ارائه‌دهندگان AI و دادهٔ مرورگر نیز در این صفحه توضیح داده شده‌اند.

این بخش صفحهٔ `security-overview` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/security-overview) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/security-overview.md) · [نسخهٔ محلی کامل](docs/upstream/docs/security-overview.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Plug-N-Meet Security and Privacy Overview

This document provides a comprehensive overview of the security architecture and data privacy principles that govern the Plug-N-Meet platform. It details not only the technical measures we take to protect your meetings from unauthorized access, but also provides a transparent look at how your data is handled, stored, and managed throughout its entire lifecycle. The system is designed from the ground up to be secure by default and to give you, the operator, full control over your data.

##### Table of Contents

1.  [API Authentication](https://www.plugnmeet.org/docs/security-overview#1-api-authentication)
2.  [User Authentication & Authorization](https://www.plugnmeet.org/docs/security-overview#2-user-authentication--authorization)
3.  [End-to-End Encryption (E2EE)](https://www.plugnmeet.org/docs/security-overview#3-end-to-end-encryption-e2ee)
    -   [Server-Generated Keys](https://www.plugnmeet.org/docs/security-overview#31-server-generated-keys)
    -   [User-Provided Keys](https://www.plugnmeet.org/docs/security-overview#32-user-provided-keys)
4.  [Transport & Data-in-Transit Security](https://www.plugnmeet.org/docs/security-overview#4-transport--data-in-transit-security)
    -   [NATS Communication](https://www.plugnmeet.org/docs/security-overview#41-nats-communication)
    -   [Media Server (LiveKit)](https://www.plugnmeet.org/docs/security-overview#42-media-server-livekit)
    -   [SIP Gateway](https://www.plugnmeet.org/docs/security-overview#43-sip-gateway)
5.  [AI Services & Data Privacy (The Insights Platform)](https://www.plugnmeet.org/docs/security-overview#5-ai-services--data-privacy-the-insights-platform)
6.  [Secure Session Flow](https://www.plugnmeet.org/docs/security-overview#6-secure-session-flow)
7.  [Browser Storage (IndexedDB)](https://www.plugnmeet.org/docs/security-overview#7-browser-storage-indexeddb)
8.  [Server-Side Data Handling & Persistence](https://www.plugnmeet.org/docs/security-overview#8-server-side-data-handling--persistence)

---

###### 1. API Authentication

All server-to-server API requests, such as creating a room or generating a join token, must be authenticated to prevent unauthorized access. This is achieved through a robust HMAC-based signature verification process. For detailed information on the API endpoints and their parameters, please refer to the [API Documentation](https://www.plugnmeet.org/docs/api/intro.md).

-   **API Key**: Every request must include a valid `API-KEY` in its headers.
-   **HMAC Signature**: A `HASH-SIGNATURE` header must be provided, which is a **HMAC-SHA256** hash of the request payload, signed with the shared `API-SECRET`.

The server validates the API key and re-computes the hash on its end. It performs a constant-time comparison against the provided signature to mitigate timing attacks. Any request with an invalid key or signature is immediately rejected.

*Reference: `server/pkg/controllers/auth.go` (`HandleAuthHeaderCheck`)*

###### 2. User Authentication & Authorization

Before a user can connect to a room's real-time infrastructure (NATS and LiveKit), their access token must be verified.

1.  The client sends its access token to the `/api/verifyToken` endpoint.
2.  The backend server validates the JWT, checks for duplicate connections, verifies the user is not on a blocklist, and ensures the room has not exceeded its participant limit.
3.  Upon successful verification, the server returns the necessary connection details, including NATS WebSocket URLs and a list of NATS subjects the user is permitted to interact with.

All subsequent API calls from the client to the backend during the session require the `Authorization` header with the user's access token.

To maintain a secure and uninterrupted session, the client application is responsible for proactively renewing the access token before it expires. This renewal process occurs securely over NATS: the client sends a renewal request via NATS, and the server responds with a new, updated token through the same secure NATS channel.

*References: `server/pkg/controllers/auth.go` (`HandleVerifyToken`), `client/src/helpers/api/plugNmeetAPI.ts`*

###### 3. End-to-End Encryption (E2EE)

Plug-N-Meet provides robust End-to-End Encryption (E2EE) for all peer-to-peer data, including media streams (audio, video, screen sharing), chat messages, whiteboard data and all notifications. This ensures that the server never has access to the unencrypted content.

###### 3.1. Server-Generated Keys

When a room is created with E2EE enabled but `enabled_self_insert_encryption_key` is `false`, the backend server securely generates a 32-byte random string to serve as the base secret. This secret is stored as part of the room's metadata and distributed to authenticated participants when they join.

-   Upon receiving this secret, the client-side application does **not** use it directly for encryption.
-   Instead, it processes the secret through the **PBKDF2 (Password-Based Key Derivation Function 2)** algorithm, using the room's random **Session ID** as the salt. This ensures that even if the same `roomId` is reused, each session will generate a unique cryptographic key.
-   This derives a strong AES-256 key that is used for all subsequent encryption and decryption operations.

*References: `protocol/utils/create_room.go` (`SetCreateRoomDefaultValues`), `protocol/proto_files/plugnmeet_create_room.proto`*

###### 3.2. User-Provided Keys

For maximum security and zero-trust, rooms can be configured with `enabled_self_insert_encryption_key` set to `true`.

-   In this mode, each user is prompted to enter a secret key or passphrase upon joining. It is the participants' responsibility to coordinate and share the same secret key, as the server has zero knowledge about the key as the plain key never leaves the user's device. If participants use different keys, they will not be able to communicate with each other in the encrypted session.
-   The client-side application processes this passphrase through the same robust **PBKDF2** algorithm, using the room's random **Session ID** as the salt. This process is intentionally slow to make brute-force attacks against the user's passphrase infeasible.
-   The strong cryptographic key derived from this process is then used for all encryption and decryption, ensuring the user's original secret is never used directly and is well-protected against attacks.

*References: `client/src/components/extra-pages/InsertE2EEKey.tsx`, `client/src/helpers/libs/cryptoMessages.ts` (`importSecretKeyFromPlainText`), `client/src/helpers/nats/ConnectNats.ts` (`createMediaServerConn`)*

###### 4. Transport & Data-in-Transit Security

###### 4.1. NATS Communication

Communication for signaling, chat, and other data messages is handled via NATS.

-   **Permissions**: When a user connects to NATS, the `nats-auth-controller` on the backend dynamically generates a set of permissions based on the user's validated token. This strictly limits which subjects a user can publish or subscribe to, enforcing a principle of the least privilege.
-   **E2EE**: If E2EE is enabled for the room, all data payloads (e.g., chat messages, whiteboard updates) are encrypted on the client side *before* being published to NATS. The server only routes the encrypted blobs of data.

*References: `server/pkg/controllers/nats_auth_controller.go`, `client/src/helpers/nats/ConnectNats.ts`*

###### 4.2. Media Server (LiveKit)

Media streams are managed by LiveKit, which has built-in support for E2EE.

-   The same cryptographic key used for NATS data is passed to the LiveKit SDK.
-   The LiveKit client SDK handles the encryption and decryption of all audio, webcam, and screen-sharing tracks directly on the user's device. The media server (SFU) only relays encrypted media packets.

*Reference: `client/src/helpers/livekit/ConnectLivekit.ts`*

###### 4.3. SIP Gateway

The telephone dial-in feature is powered by the **[livekit/sip](https://github.com/livekit/sip)** project. This is an advanced feature that requires the server administrator to configure it with an external SIP trunking provider to bridge calls between the telephone network (PSTN) and your WebRTC session.

-   **Administrator Responsibility:** You are responsible for the setup and security of this integration. For detailed setup instructions, please refer to the official [livekit/sip GitHub repository](https://github.com/livekit/sip).
-   **Transport Encryption (SRTP):** The `livekit/sip` service supports `SIPMediaEncryption` (SRTP), which encrypts the audio stream between the SIP gateway and your provider. It is your responsibility to enable this and to ensure your provider supports it.
-   **E2EE Incompatibility:** The SIP gateway must process unencrypted audio to bridge the call between the PSTN and WebRTC networks. Because of this, it is **fundamentally incompatible with End-to-End Encryption**. The SIP dial-in feature is automatically disabled for any room where E2EE is active.

###### 5. AI Services & Data Privacy (The Insights Platform)

The Plug-N-Meet Insights Platform integrates with third-party AI providers (such as Microsoft Azure and Google) to offer advanced features. As the operator of the platform, it is your responsibility to understand how this data is handled and to update your own privacy policy accordingly.

###### 5.1. Audio-Based AI Features

Features that require understanding spoken words, such as **Speech-to-Text**, **Spoken Translations**, and **AI Meeting Summaries**, work by processing the meeting's audio.

*   **How it Works:** When these features are enabled, an AI agent joins the session. This agent captures the mixed audio stream of the meeting and sends it directly to the configured AI provider for processing.
*   **Data Flow:** The audio data is streamed from your Plug-N-Meet server to the third-party AI provider's servers.

###### 5.2. Text-Based AI Features

Text-based features have different data flow models depending on their purpose.

*   **Chat Translation:** This is a **broadcast** feature. When a user sends a chat message, the text is sent to the AI provider for translation. The translated text is then received by the Plug-N-Meet server and broadcast to the relevant participants.
*   **AI Chat Assistant:** This is a **private, one-to-one** feature. When a user sends a message to the AI Assistant, that prompt is sent to the AI provider. The response is sent back **only to that specific user** and is not broadcast to anyone else in the meeting.

###### 5.3. The Impact of End-to-End Encryption (E2EE)

Your privacy and security are paramount. The interaction between AI features and E2EE is designed to be secure by default.

*   **When E2EE is enabled with a self-provided key (`enabled_self_insert_encryption_key: true`), all audio-based AI features are automatically disabled.** This is a deliberate security design. The AI agent cannot decrypt the audio stream, and therefore cannot function.
*   **Text-based features like Chat Translation and the AI Chat Assistant will continue to function.** This is because the text is sent to the AI provider *before* it is encrypted for transport over NATS.

**Disclaimer for Operators:** You are the data controller. It is your responsibility to inform your users that when AI features are enabled, their audio and/or text data will be processed by third-party AI providers. You must ensure your terms of service and privacy policy are updated to reflect this data flow.

###### 6. Secure Session Flow

The end-to-end connection process is designed with security at each step. The exact flow can vary based on the room's configuration (e.g., whether E2EE is enabled).

1.  **Token Verification**: The client verifies its access token with the backend.
2.  **NATS Connection**: The client connects to NATS using the credentials from step 1. The NATS auth service grants fine-grained permissions for the session.
3.  **Initial Data Fetch**: The client requests initial room data over the secure NATS connection.
4.  **(Conditional) E2EE Key Import**: If End-to-End Encryption is enabled for the room, the client imports the E2EE key. This key is either received from the server or derived from the user's manual input, depending on the room's E2EE settings.
5.  **Media Connection**: The client connects to the LiveKit media server. If E2EE is active, it enables encryption with the imported key.
6.  **Secure Communication**: All subsequent data and media are now handled according to the room's security configuration, with E2EE applied if it was enabled.

###### 7. Browser Storage (IndexedDB)

To enhance user experience and ensure session continuity, the application utilizes the browser's built-in IndexedDB storage. This is used for purely functional purposes, such as restoring a user's session if they accidentally refresh their page.

To achieve this, the application temporarily stores ephemeral session data on the end-user's local device. The categories of data stored may include, but are not limited to:

*   **User & Session Preferences:** To remember settings like a chosen subtitle language.
*   **Session Content & History:** To restore chat messages, speech-to-text transcripts, and the state of the whiteboard.
*   **Performance Caches:** Such as cached images to improve loading times.

###### Data Lifecycle and Encryption

The data stored in IndexedDB is designed to be ephemeral and is handled as follows:

*   **Encryption:** The data stored in IndexedDB is **not encrypted** by the application. This is a deliberate design choice, as the data resides on the end-user's own device and is protected by the security measures of the user's operating system and browser profile. The primary security focus is on encrypting data *in transit* and ensuring it is not persisted on any server.
*   **Normal Session End:** When a user properly ends their session (e.g., by clicking "End meeting" or "Leave Meeting"), all stored data for that session is **immediately and permanently deleted** from the browser.
*   **Abnormal Session End (e.g., closing the browser tab):** If a session is not ended properly, the data remains in the browser's IndexedDB. However, upon the next application startup, a cleanup process runs. This process automatically identifies and deletes any stored data from previous sessions that has exceeded a predefined maximum age (e.g., several hours). This maximum age is configurable and subject to change in future versions to optimize performance and privacy.

**Key Points:**

*   This data is stored **only on the end-user's browser** and is never transmitted elsewhere for storage purposes.
*   The storage is essential for the application's expected functionality and is not used for tracking or analytics.

**Disclaimer for Operators:** As the person or organization deploying this software, you are responsible for creating and maintaining your own Privacy Policy. You should use this information to ensure your policy is transparent and compliant with any applicable data protection regulations (e.g., GDPR, CCPA).

###### 8. Server-Side Data Handling & Persistence

To provide a complete picture of the data lifecycle, this section outlines how data is handled and stored on the server side. The architecture is designed to separate volatile real-time data from long-term persistent data, with a strong emphasis on user-controlled data retention.

###### 8.1. Real-Time Session State (Redis or NATS KV)

During an active meeting, the server needs to manage the real-time state of the room, such as the list of participants, their mute status, and other immediate metadata. This is handled by a high-performance, in-memory key-value store (either Redis or NATS KV).

-   **Purpose:** Fast, real-time coordination of an active session.
-   **Lifecycle:** This data is volatile and is tied to the life of the session. It is automatically cleared when the session ends.

###### 8.2. Temporary Session Cache for Data Sync

To improve session reliability, the server temporarily caches certain real-time collaboration data, such as whiteboard/notepad snapshots.

-   **Purpose:** This cache acts as a fallback to ensure a smooth user experience. For example, if a user joins late or switches between whiteboard pages, they can quickly fetch the latest content from the server instead of depending on other participants to be online and send it to them.
-   **E2EE Compliance:** This feature is designed to be fully compatible with End-to-End Encryption. When E2EE is active in a room, the data stored in this cache is the **encrypted ciphertext**. The server only holds an unreadable blob of data and has no ability to decrypt it, preserving the zero-knowledge privacy model.
-   **Lifecycle:** This data is highly volatile. It is stored in the same in-memory store as the real-time session state and is **automatically and permanently deleted** the moment the session ends. It is never written to a persistent database or the filesystem.

###### 8.3. Persistent Database (MariaDB)

For historical reference and administrative purposes, a small subset of non-sensitive information is stored in a persistent relational database (MariaDB).

-   **Purpose:** Long-term record-keeping of meeting occurrences.
-   **Data Stored:** This typically includes basic room information such as `roomId`, `title`, creation/end time etc.

###### 8.4. Optional Analytics Data

PlugNmeet provides the option to persist detailed analytics for a session to help administrators understand usage patterns. This feature is governed by a setting that provides administrators with control over data retention.

-   **User Control:** The decision to **persist** analytics data is made by the administrator on a per-room basis via the `enable_analytics` flag, adhering to the principle of privacy by default.
-   **Data Collected (Metadata, Not Content):** To protect user privacy, PlugNmeet **does not store the raw content** of user interactions in its analytics. Instead, it aggregates metadata and event counters. For example, the system records:
    -   The *number of times* a user sent a chat message, not the chat messages themselves.
    -   The *number of times* a user drew on the whiteboard, not the content of the drawings.
    -   The *total duration* a user spoke, not the audio of what they said.
    -   Other engagement events like join/leave times, files uploaded, and poll votes.
-   **Storage:** If persisted, this aggregated metadata is **always stored as a JSON file** on the server's filesystem. A reference to the filename and its associated room is then stored in a dedicated table in the persistent database for easy retrieval.
-   **Lifecycle:** Analytics data is aggregated in-memory during the course of an active session. Upon the conclusion of the session, the system checks the room's `enable_analytics` setting.
    -   If `true`, the aggregated metadata is written to a JSON file, and its reference is saved to the database.
    -   If `false`, the in-memory data is immediately discarded and is not persisted in any form.

###### 8.5. Cloud Recordings

When cloud recording is enabled for a session, the resulting media file (MP4) is stored on the server.

-   **Purpose:** To provide a persistent, shareable record of a meeting.
-   **Storage:** The MP4 files are stored in a configurable directory on the server's filesystem. To allow for easy management and retrieval via the API, a reference to each recording—including its `record_id`, associated `roomId`, `file_path`, `file_size`, and creation timestamp—is stored in a dedicated table in the persistent database.
-   **User Control & Lifecycle:** The management of these recordings is entirely controlled by the administrator via the API. Recordings are retained on the server indefinitely until they are explicitly deleted using the `/recording/delete` API call. This gives the administrator full control over the data retention lifecycle.
-   **E2EE Incompatibility:** Server-side recording is **fundamentally incompatible** with zero-trust End-to-End Encryption. Therefore, cloud recording is automatically disabled and cannot be initiated if the room is configured with `enabled_self_insert_encryption_key: true`. This is because the server has no access to the unencrypted media streams required to create the recording, which is the core guarantee of the E2EE model.

###### 8.6. AI Service Artifacts

When AI features like transcription or summarization are used, they generate valuable data artifacts. Similar to cloud recordings and analytics, Plug-N-Meet gives the administrator full control over this data.

-   **Purpose:** To provide a persistent, retrievable record of AI-generated meeting intelligence.
-   **Storage:** The generated files (e.g., VTT transcription files, text summaries) are stored in a configurable directory on the server's filesystem. To allow for easy management and retrieval, **metadata** about these artifacts—such as the associated `roomId`, `file_path`, AI token usage, and speech duration—is stored in the `pnm_room_artifacts` table in the persistent database.
-   **User Control & Lifecycle:** The management of these artifacts is entirely controlled by the administrator via the API. You can:
    -   **Download** the transcription or summary files.
    -   **Delete the Artifact File:** The API allows you to delete the generated files from the server's filesystem to save storage space.
    -   **Retained Metadata:** The associated metadata in the `pnm_room_artifacts` table is linked to the meeting's history and is **retained** for your records, even after the file is deleted. This ensures you always have a permanent record for auditing and analyzing AI service usage and costs.

---

This layered approach to security and data handling does more than just protect against common threats; it provides a transparent and flexible framework that puts you, the operator, in control. By combining robust security measures with a "**_privacy by design_**" philosophy, Plug-N-Meet gives you the tools and the transparency needed to meet your own privacy and compliance obligations confidently.

<a id="page-39"></a>

## سازگاری و مهاجرت از BBB

راهنما امکان استفاده از مسیر سازگار BBB و تغییر تنظیمات اتصال برنامه‌های موجود را توضیح می‌دهد. برای طرح آی‌نشست، ارائه‌دهنده باید روی جلسه ذخیره شود تا ساخت و پایان و ضبط هر جلسه از همان مسیر انجام شود. قرارداد اختصاصی فایل‌های ضبط و بازه‌های گفتار آی‌نشست باید جداگانه تأمین و آزموده شود. تبدیل URL به‌تنهایی تأییدکنندهٔ برابری خروجی‌های تحلیلی نیست.

این بخش صفحهٔ `tutorials/migration-from-bbb` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/tutorials/migration-from-bbb) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/tutorials/migration-from-bbb.md) · [نسخهٔ محلی کامل](docs/upstream/docs/tutorials/migration-from-bbb.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Migrating from BigBlueButton: A Modern Alternative

Are you looking for a more modern, lightweight, and scalable alternative to BigBlueButton? Plug-N-Meet offers a next-generation architecture with a focus on performance and simplicity, and we have made the migration process incredibly straightforward.

This guide will walk you through the benefits of switching and show you how to migrate your existing BBB-integrated application in just a few minutes, without changing any of your front-end code.

---

##### Why Consider Plug-N-Meet?

BigBlueButton is a foundational open-source project that has served the community for years. Plug-N-Meet was born from a deep respect for the BBB ecosystem and was designed to address some of its inherent architectural challenges, offering a different approach for those who need more flexibility and scalability.

For a detailed, head-to-head comparison of the two platforms, we recommend reading our full breakdown:
**[BigBlueButton vs. Plug-N-Meet: A Modern Alternative](https://www.plugnmeet.org/blog/bigbluebutton-vs-plugnmeet-alternative)**

Here are the key advantages of the Plug-N-Meet architecture:

-   **Built on a World-Class Foundation:** Plug-N-Meet is powered by **LiveKit** as its media server and **NATS** for real-time messaging (chat, whiteboard sync, etc.). These are extremely fast, scalable, and modern components that allow Plug-N-Meet to handle a massive user load with fewer server resources. See measured numbers in our [Benchmarks](https://www.plugnmeet.org/docs/benchmark): 300–500 users in one room ran clean on a single 4-vCore server.
-   **Simplified & Lightweight Stack:** By leveraging these components, Plug-N-Meet avoids the complex web of dependencies found in BigBlueButton (like FreeSWITCH and Kurento). This results in a significantly smaller footprint, easier installation, and simpler maintenance.
-   **Simplified, Single-File Recordings:** Unlike BigBlueButton's complex post-processing workflow, plugNmeet generates a standard MP4 file for each recording. This makes your recordings instantly portable, easier to manage, and faster to access.
-   **Designed for Distributed Scaling:** Plug-N-Meet is not built as a single, large application. Each core component—LiveKit (media), NATS (messaging), the recorder, and the plugnmeet server—is designed to be scaled out independently. This allows you to add resources precisely where you need them, for example, scaling your recording infrastructure without impacting live meeting performance.
-   **Effortless Customization:** The platform is designed for deep and easy customization. The frontend client is a modern, self-contained application, and the backend is a single binary. This clean separation allows developers to rapidly change the UI, add features, or integrate custom workflows without navigating the complex, multi-service environment found in BigBlueButton.
-   **Seamless API Compatibility:** This is our commitment to the ecosystem. Plug-N-Meet includes a BBB-compatible API layer. This means your existing application—whether it's Greenlight, a Moodle plugin, or a custom-built solution—will work with Plug-N-Meet without requiring any front-end code changes.

---

##### The 2-Minute Migration Process

Because of our BBB-compatible API, migrating your existing application is as simple as changing the API endpoint and credentials in your configuration.

In your application's settings, locate your BigBlueButton API configuration and update the following two values:

1.  **API Endpoint URL:** Change the URL to point to your PlugNmeet server, followed by `/bigbluebutton/`.
2.  **API Secret:** Change the secret to match your PlugNmeet API Secret.

That's it. Your `API Key` from PlugNmeet will be passed as the `user` in the URL.

###### Example Configuration

```
URL: https://[YOUR_PLUGNMEET_HOST]/[YOUR_PLUGNMEET_API_KEY]/bigbluebutton/
Secret: [YOUR_PLUGNMEET_API_SECRET]
```

For example, using the public demo server, the configuration would be:

```
URL: https://demo.plugnmeet.com/plugnmeet/bigbluebutton/
Secret: EXAMPLE_API_SECRET
```

---

##### After Migration: Unlocking the Full Feature Set

To provide this seamless migration, Plug-N-Meet uses a powerful, BBB-compatible API layer. This ensures that all of BigBlueButton's core features that your existing application relies on will work out-of-the-box.

However, this compatibility layer is designed as a bridge. To experience the full power of Plug-N-Meet—including our most advanced features, superior performance, and effortless customization options—we **strongly recommend switching to our dedicated plugins or SDKs** when you are ready.

Using the native integrations will provide the most stable, feature-rich, and future-proof experience for your users.

-   **[WordPress Plugin](https://www.plugnmeet.org/docs/user-guide/wordPress-integration.md)**
-   **[Moodle Plugin](https://www.plugnmeet.org/docs/user-guide/moodle-integration.md)**
-   **[Joomla Plugin](https://www.plugnmeet.org/docs/user-guide/joomla-integration.md)**
-   **[PHP & JavaScript SDKs](https://www.plugnmeet.org/docs/tutorials/quick_php.md)** for custom applications.

<a id="page-40"></a>

## نمونهٔ PHP برای ایجاد و ورود

الگوی نمونه ابتدا فعال‌بودن اتاق را بررسی می‌کند، در صورت نیاز اتاق می‌سازد و سپس برای کاربر مجاز توکن ورود می‌گیرد. secret در backend می‌ماند و مرورگر فقط توکن کوتاه‌عمر خود را دریافت می‌کند. در برنامهٔ واقعی رقابت دو درخواست ساخت و خطای هر گام باید مدیریت شود. نمونهٔ SDK و redirect در پیوست موجود است.

این بخش صفحهٔ `tutorials/quick_php` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/tutorials/quick_php) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/tutorials/quick_php.md) · [نسخهٔ محلی کامل](docs/upstream/docs/tutorials/quick_php.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### PHP Quick Start: Create and Join a Meeting

This tutorial provides a complete, step-by-step example of how to use the `plugNmeet-sdk-php` library to create a new meeting room and generate a secure join token for a user.

The logic follows a common workflow: check if a room exists, create it if it doesn't, and then generate a URL for a user to join.

<img src="/img/tutorials/quick_join_flow.png" width="400" alt="Quick Join Logic Flow" loading="lazy"/>

---

##### Prerequisites

Before you begin, ensure you have the following:

- A working PlugNmeet server with your API Key and Secret.
- The latest version of the [plugNmeet-sdk-php](https://github.com/mynaparrot/plugNmeet-sdk-php/releases) library downloaded.

---

##### Step 1: Setup & Configuration

First, create a PHP file (e.g., `quickJoin.php`) and include the `plugNmeetConnect.php` class from the SDK. Then, create a configuration object with your server details.

```php
<?php
require __DIR__ . "/plugNmeetConnect.php";

// Step 1: Configuration
$config = new stdClass();
$config->plugnmeet_server_url = "http://localhost:8080"; // Your server URL
$config->plugnmeet_api_key = "plugnmeet"; // Your API Key
$config->plugnmeet_secret = "EXAMPLE_API_SECRET"; // Your API Secret

$connect = new plugNmeetConnect($config);
```

---

##### Step 2: Define Room & User Parameters

Next, define the basic parameters for the room you want to create and the user who will be joining. The `$roomMetadata` array allows you to customize every aspect of the room's features, from enabling webcams to setting default lock permissions.

```php
// Step 2: Room and User definitions
$roomId = "room01"; // Must be unique. You can also use $connect->getUUID();
$max_participants = 0; // 0 = no limit
$user_full_name = "Your Name";
$userId = "your-unique-user-id"; // Must be unique for each user.
$isAdmin = true;

// Define all the features for this specific room.
$roomMetadata = array(
    "room_features" => array(
        "allow_webcams" => true,
        "mute_on_start" => false,
        "allow_screen_share" => true,
        "allow_rtmp" => true,
        "allow_view_other_webcams" => true,
        "allow_view_other_users_list" => true,
        "admin_only_webcams" => false,
        "enable_analytics" => true,
        "room_duration" => 0, // in minutes. 0 = no limit/unlimited
        "allow_virtual_bg" => true,
        "allow_raise_hand" => true,
        // if false then you'll need to provide unique user id
        "auto_gen_user_id" => true,
        "recording_features" => array(
            "is_allow" => true,
            "is_allow_cloud" => true,
            "is_allow_local" => true,
            "enable_auto_cloud_recording" => false
        ),
        "chat_features" => array(
            "is_allow" => true,
            "is_allow_file_upload" => true
        ),
        "shared_note_pad_features" => array(
            "is_allow" => true
        ),
        "whiteboard_features" => array(
            "is_allow" => true,
            //"preload_file" => "https://mydomain.com/text_book.pdf"
        ),
        "external_media_player_features" => array(
            "is_allow" => true
        ),
        "waiting_room_features" => array(
            "is_active" => false,
        ),
        "breakout_room_features" => array(
            "is_allow" => true,
            "allowed_number_rooms" => 2
        ),
        "display_external_link_features" => array(
            "is_allow" => true,
        ),
        "ingress_features" => array(
            "is_allow" => true,
        ),
        "polls_features" => array(
            "is_allow" => true,
        ),
        "insights_features" => array(
            "is_allow" => true,
            "transcription_features" => array(
                "is_allow" => true,
                "is_allow_translation" => true,
                "is_allow_speech_synthesis" => true,
            ),
            "chat_translation_features" => array(
                "is_allow" => true,
            ),
            "ai_features" => array(
                "is_allow" => true,
                "ai_text_chat_features" => array(
                    "is_allow" => true,
                ),
                "meeting_summarization_features" => array(
                    "is_allow" => true,
                )
            )
        ),
        "end_to_end_encryption_features" => array(
            "is_enabled" => false,
            "enabled_self_insert_encryption_key" => false,
            "included_chat_messages" => false,
            // this may use more CPU for the user end.
            // do not enable it unless really necessary
            "included_whiteboard" => false,
        ),
    ),
    "default_lock_settings" => array(
        "lock_microphone" => false,
        "lock_webcam" => false,
        "lock_screen_sharing" => true,
        "lock_whiteboard" => true,
        "lock_shared_notepad" => true,
        "lock_chat" => false,
        "lock_chat_send_message" => false,
        "lock_chat_file_share" => false,
        "lock_private_chat" => false // user can always send private message to moderator
    ),
    // copyright_conf will only work if server config has been
    // set true for `allow_override` otherwise this will ignore
    "copyright_conf" => array(
        "display" => true,
        "text" => "Powered by <a href=\"https://www.plugnmeet.org\" target=\"_blank\">plugNmeet</a>"
    )
);
```

---

##### Step 3: The Logic Flow - Check, Create, Join

The following code blocks implement the core logic.

###### 3.1 Check if the Room is Active

First, we call the API to see if a room with the specified `$roomId` already exists and is active.

```php
$isRoomActive = false;
$output = new stdClass();
$output->status = false;

try {
    $res = $connect->isRoomActive($roomId);
    if (!$res->getStatus()) {
        $output->msg = $res->getMsg();
    } else {
        $isRoomActive = $res->getIsActive();
        $output->status = true;
    }

} catch (Exception $e) {
    $output->msg = $e->getMessage();
}
```

###### 3.2 If Not, Create the Room

If the room is not active, we proceed to create it using the parameters we defined earlier.

```php
if (!$isRoomActive && $output->status) {
    try {
        $create = $connect->createRoom($roomId, "Test room", $roomMetadata, "Welcome to room", "", "", $max_participants);

        $isRoomActive = $create->getStatus();
        $output->status = $create->getStatus();
        $output->msg = $create->getMsg();
    } catch (Exception $e) {
        $output->msg = $e->getMessage();
    }
}
```

###### 3.3 Generate the Join Token

Once we have confirmed the room is active (either because it already existed or we just created it), we can generate a secure, single-use join token for our user.

```php
if ($isRoomActive && $output->status) {
    try {
        $join = $connect->getJoinToken($roomId, $user_full_name, $userId, $isAdmin);

        if ($join->getStatus()) {
            $output->token = "<br>" . $join->getToken();
            $output->url = "<br>" . $config->plugnmeet_server_url . "?access_token=" . $join->getToken();
            // or you can set cookie name `pnm_access_token` with that token & redirect
        }
        $output->status = $join->getStatus();
        $output->msg = $join->getMsg();
    } catch (Exception $e) {
        $output->msg = $e->getMessage();
    }
}
```

---

##### Step 4: Redirecting the User

After successfully generating the join token, the complete meeting URL will be in the `$output->url` variable. You can now redirect your user to this URL to join the meeting.

```php
if ($output->status) {
    header("Location: " . $output->url);
    exit;
} else {
    echo $output->msg;
}
```

---

##### Next Steps

This quick start guide uses the simplest method to get a user into a room. For a more deeply integrated experience, you can build a custom client interface using the [getClientFiles()](https://www.plugnmeet.org/docs/api/get-client-files) API method. This allows you to host the client on your own page without using an iframe.

See the [conference.php](https://github.com/mynaparrot/plugNmeet-sdk-php/blob/main/examples/conference.php) file in our PHP-SDK for a complete example.

<a id="page-41"></a>

## راهنمای شرکت‌کننده

شرکت‌کننده پس از ورود مجوز میکروفن و دوربین مرورگر را می‌دهد و می‌تواند صوت و تصویر خود را با توجه به قفل‌های مدیر کنترل کند. گفتگو، فهرست کاربران، پاسخ به نظرسنجی، مشاهدهٔ محتوا و جابه‌جایی اتاق فرعی در این راهنما پوشش داده شده‌اند. امکانات ترجمه، دستیار و رونویسی تابع تنظیمات اتاق هستند. خروج شرکت‌کننده باید در گزارش حضور با شناسهٔ همان شخص ثبت شود.

این بخش صفحهٔ `user-guide/attendee` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/user-guide/attendee) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/user-guide/attendee.md) · [نسخهٔ محلی کامل](docs/upstream/docs/user-guide/attendee.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Attendee Guide: How to Join and Participate in a Meeting

Welcome to the Attendee Guide! This document will walk you through everything you need to know to successfully join and participate in a Plug-N-Meet video conference.

---

##### First Steps: Joining a Meeting

When you first enter a PlugNmeet room, you will be prompted to choose how you want to join the audio session. This is your most important first step.

![Onboarding-Desktop.jpg](https://www.plugnmeet.org/img/moderator/Onboarding-Desktop.jpg)

-   **Enable Microphone and Camera:** Choose this option to speak and share your video in the meeting. Your browser will ask for permission to access your microphone and camera.
-   **Continue as a listener:** Choose this option to join without a microphone or camera. You will be able to hear everything but will not be able to speak or share video.

**Important:** This choice is not permanent. If you decide to **Continue as a listener**, you can easily activate your microphone and camera later from the main control bar (if the moderator allows it).

---

##### Section 1: The Attendee Control Bar

These are the basic controls you will use during the meeting, located at the bottom of the screen.

-   **Microphone:** Click this icon to mute or unmute yourself.
-   **Webcam:** Click this icon to start or stop sharing your webcam. You can choose a virtual background before you share.
-   **Raise Hand:** Click this icon to signal to the moderator that you have a question or wish to speak.

![attendee-control-bar.webp](https://www.plugnmeet.org/img/moderator/attendee-control-bar.webp)

---

##### Section 2: Participating in the Meeting

These features allow you to interact with the moderator and other participants.

###### Using the Chat

-   **Public Chat:** Click the **Chat** icon to open the main chat panel where you can send messages to everyone.
-   **Private Chat:** To send a private message, open the **Participants** list, click the three-dot menu next to a user's name, and select **Send private message**.

![attendee-Chat.webp](https://www.plugnmeet.org/img/moderator/attendee-Chat.webp)

###### Viewing the Participant List

-   Click the **Participants** icon to see a list of everyone currently in the meeting.

###### Responding to Polls

-   If the moderator starts a poll, you will see a notification. You can open the **Polls** tab in the participant panel to cast your vote.

![attendee-polls.webp](https://www.plugnmeet.org/img/moderator/attendee-polls.webp)

###### Transcription & Translation
- **Start Transcription & Translation:** Click the **Transcription & Translation** icon in the main footer control bar to start the service (if the moderator has enabled it for you).

![attendee-transcriptions.webp](https://www.plugnmeet.org/img/moderator/attendee-transcriptions.webp)

- **History:** You can view **Live Transcription & Translation** in real-time. Click the **History** icon to see the full transcript history.

![attendee-transcriptions-history.webp](https://www.plugnmeet.org/img/moderator/attendee-transcriptions-history.webp)

**Important:** You must **Enable Microphone** to use **Transcription & Translation**.

###### Using AI Assistant

- **Enable/Disable:** Click the **AI Assistant** icon in the main footer control bar to start the AI Assistant (if the moderator has enabled it for you).

![attendee-ai-tools.webp](https://www.plugnmeet.org/img/moderator/attendee-ai-tools.webp)

---

##### Section 3: Viewing Content

The moderator can share various types of content during the meeting.

-   **Viewing a Shared Screen:** If the moderator shares their screen, it will automatically appear in the main area of the interface.
-   **Interacting with the Whiteboard:** If the moderator activates the whiteboard, you will be able to see what they are drawing. They may also give you permission to draw on the whiteboard yourself.
-   **Viewing a Shared Video:** If the moderator shares a video, it will appear on your screen with playback controls.

---

##### Section 4: Breakout Rooms

The moderator can split the main session into smaller Breakout Rooms for group activities.

-   If you are assigned to a room, you will see a notification inviting you to **Join**.
-   You will be taken to a new, smaller meeting room. When the breakout session is over, you will be returned to the main room automatically.

![attendee-breakout-room.webp](https://www.plugnmeet.org/img/moderator/attendee-breakout-room.webp)

---


##### Section 5: Leave Meeting

- **Leave Meeting:** If you select **Leave Meeting** from the top-right menu or the bottom-right button, you will leave the meeting, but the session will remain active for other participants.

![attendee-end-meeting.webp](https://www.plugnmeet.org/img/moderator/attendee-end-meeting.webp)

<a id="page-42"></a>

## اتصال Joomla

پس از نصب extension و ثبت مشخصات API، اتاق‌ها در دسته‌بندی قرار می‌گیرند و مجوزهای دسترسی و نمایش در سایت تعیین می‌شوند. مدیریت ضبط و artifact، سفارشی‌سازی ظاهر و پیش‌فرض کیفیت رسانه در همین راهنما آمده است. اعتبار API باید در سمت سرور سایت میزبان نگهداری شود.

این بخش صفحهٔ `user-guide/joomla-integration` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/user-guide/joomla-integration) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/user-guide/joomla-integration.md) · [نسخهٔ محلی کامل](docs/upstream/docs/user-guide/joomla-integration.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### How to Add Video Conferencing to Joomla with the Plug-N-Meet Extension

This guide covers the plugNmeet component for Joomla. It brings the same full feature set as the WordPress and Moodle plugins to Joomla: full room control, recordings with merge, artifacts with analytics export, and Joomla permission-based access — all inside your Joomla site.

---

##### 1. Key Features of the Joomla Integration

- **Full Room Control:** Create rooms with all plugNmeet features — waiting room, breakout rooms, polls, whiteboard, E2EE, SIP dial-in, and Insights AI (transcription, chat translation, AI summary).
- **Backend Management:** Manage rooms, recordings, and artifacts without leaving **Components > Plug N Meet**.
- **Recordings with Merge:** List, play, download, delete, and **merge multiple segments** into one MP4 from the Recordings screen.
- **Artifacts & Analytics:** View analytics reports, transcriptions, and AI summaries per room. Download analytics as Excel.
- **Joomla-Native Publishing:** Rooms are Joomla articles-like items with categories, plus **Single room** and **List rooms** menu item types.
- **Permission-Based Access:** Control who can join without password, join as moderator, and view/play/download/delete recordings per room via Joomla ACL.

##### 2. Installation & Initial Setup

###### Step 2.1: Install the Extension

- From your Joomla administrator panel, navigate to **System > Install > Extensions**.
- Select the **Install from Web** tab and search for **plugnmeet**.
- Find the official Plug-N-Meet extension and click **Install**.

![Joomla-install.png](https://www.plugnmeet.org/img/Joomla/Joomla-install.png)

- After a moment, you should see a confirmation that the installation was successful.

![Joomla-success.png](https://www.plugnmeet.org/img/Joomla/Joomla-success.png)

Alternatively, download the package from the [plugNMeet-Joomla GitHub releases](https://github.com/mynaparrot/plugNMeet-Joomla) and install via **Upload Package File**.

###### Step 2.2: Connect to Your Server (API Settings)

This is the most important step. You must connect the extension to a working plugNmeet server.

- From your Joomla administrator panel, navigate to **System > Global Configuration**.
- In the left-hand menu, select **plugNmeet**.

![Joomla-plugin-config.png](https://www.plugnmeet.org/img/Joomla/Joomla-plugin-config.png)

- You will see: **Server URL**, **API Key**, **API Secret**, and **Client load** (`Remote` or `Redirect`).
- After installation, these fields are pre-filled with demo credentials. **You must replace these with the API details from your own self-hosted or plugNmeet Cloud server.** The demo server is for temporary testing only and has limitations.
- **Client load:** `Remote` (default) keeps the meeting embedded inside your Joomla site. `Redirect` sends users straight to your plugNmeet server interface (`server-url/?access_token=...`) and takes them off your site layout. Joomla does not offer a Local client option like WordPress.

![Joomla-plugnmeet-config.png](https://www.plugnmeet.org/img/Joomla/Joomla-plugnmeet-config.png)

- After entering your own credentials, click **Save**.
- Global defaults also live here: **Copyright** (footer branding), **Design** (logo, colors, background, custom CSS), and **Advanced** (simulcast, dynacast, resolutions, audio preset, codec). Rooms inherit these unless overridden per room.

---

##### 3. Creating & Publishing Your First Meeting

###### Step 3.1: Organize with Categories (Optional)

Rooms can be grouped with Joomla categories:

- Go to **Components > Plug N Meet > Rooms Categories** and create categories (e.g., Classes, Meetings).
- Assign a category in the **Category** field when creating a room.
- The frontend **List rooms** menu type can filter by one category.

###### Step 3.2: Create a Meeting Room

- From the Joomla administrator panel, go to **Components > Plug N Meet > Manage Rooms**.
- Click **+ New** to create a new meeting room.

![Joomla-12.png](https://www.plugnmeet.org/img/Joomla/Joomla-12.png)

The room form has the following tabs:

**Basic**
- **Title, Category & Description:** Title shown on the site; description shown on the room page.
- **Passwords:** Separate **Moderator** and **Attendee** passwords. They must be different.
- **Welcome Message:** Shown in chat on join.
- **Max Participants:** `0` for unlimited.

**Room Features**
Core in-room options: webcams, mute on start, screen share, view other webcams/users list, admin-only webcams, analytics, virtual background, raise hand, reactions, auto-generate user ID, and room duration (minutes, `0` = no limit). **Moderator join first:** when enabled, attendees see "meeting has not started yet" until a moderator joins — different from Waiting Room, where users can already wait inside.

**Recording Features**
Master switch plus cloud recording, auto-start cloud recording, and local recording.

**Chat / Shared Notepad / Whiteboard / External Media Player / Display External Link**
Toggle each tool. Whiteboard supports a **Preload File** URL (PDF/presentation).

**Waiting Room**
Enable to screen participants, plus a custom waiting room message.

**Breakout Room**
Enable and set max rooms (up to 10).

**Ingress / External Broadcasting**
Allow RTMP ingress and external broadcasting (with separate RTMP toggle).

**End-to-End Encryption**
Enable E2EE plus options for self-insert key, chat, and whiteboard inclusion. Note: E2EE disables some server-side features such as cloud recording and transcription.

**Polls / SIP Dial-In**
Enable polls. SIP dial-in adds options to enable on create and hide phone number (requires SIP configured on the server).

**Insights Features**
Master switch plus transcription (with translation and speech synthesis), chat translation, AI text chat, and AI meeting summarization. These require Insights/AI enabled on your plugNmeet server.

**Default Lock Settings**
Initial locked state for microphone, webcam, screen sharing, whiteboard, shared notepad, chat, private chat, chat send message, chat file share, and reactions. Moderators can unlock live.

**Design Customisation**
Per-room logo, background image, primary/secondary/background/header/footer/side-panel colors, and custom CSS URL (HTTPS). Empty = inherit the global **Design** values from Global Configuration.

**Advanced Settings**
Per-room override of the global **Advanced** defaults: simulcast, dynacast, stop mic track on mute, default webcam resolution (`h90`–`h2160`), screen share resolution, audio preset (`telephone`–`musicHighQualityStereo`), and video codec (`vp8`, `h264`, `vp9`, `av1`). Leave a field empty to keep the global value.

Click **Save** when done.

###### Step 3.3: Permissions

Click the **Permissions** tab when editing a room to control Joomla group access per room:

- **Join passwordless** (`join.passwordless`): allowed groups can join without entering a password (name only).
- **Join as admin** (`join.admin`): allowed groups join as moderator even with the attendee password or no password.
- **Recording view / play / download / delete** (`recording.*`): granular control over the frontend recordings list. Guests should never get delete permission.

![room_permission.png](https://www.plugnmeet.org/img/Joomla/room_permission.png)

Use **Components > Plug N Meet** global **Options > Permissions** for site-wide defaults.

###### Step 3.4: Display the Room on Your Site

You have two menu item types under **Plug N Meet**:

**Single room** — for one meeting:

- Navigate to **Menus > Main Menu** and click **+ Add New Menu Item**.
- In **Menu Item Type**, click **Select**, then choose **Plug N Meet > Single room**.

![Joomla-25.png](https://www.plugnmeet.org/img/Joomla/Joomla-25.png)

- In the **Select a room** dropdown, choose the room you just created.
- Give the menu item a **Title** (e.g., "Live Class") and click **Save**.

![Joomla-26.png](https://www.plugnmeet.org/img/Joomla/Joomla-26.png)

**List rooms** — for a directory of rooms:

- Choose **Plug N Meet > List rooms** as the menu item type.
- Optionally set **Choose a Category** to filter to one category.

---

##### 4. The Attendee Experience

When a user clicks the menu item, they see the room title, description, and a join form.

- They enter their name and the room password, then click **Join**. With **Client load = Remote**, the meeting opens in a new tab; closing it returns via `?returned=true`. With **Client load = Redirect**, they leave Joomla and join directly on your plugNmeet server URL.
- Groups with **Join passwordless** permission see no password field.
- If **Moderator join first** is enabled and no moderator has started the room, attendees see "The meeting has not started yet" and cannot enter.
- If the user has **Recording view** permission, a recordings table appears below the join form with **Play / Download / Delete** buttons according to their permissions.

![Joomla-28.png](https://www.plugnmeet.org/img/Joomla/Joomla-28.png)

##### 5. Post-Session Management (Admin)

All assets are managed from **Components > Plug N Meet** without leaving Joomla.

###### Managing Recordings

Go to **Components > Plug N Meet > Recordings**:

1. Select a room from the dropdown and click **Show**.
2. Use checkboxes to select segments and click **Merge** to combine them into one MP4 (confirmation shows the selected order).
3. Use per-row **Play / Download / Delete** actions.

###### Accessing Artifacts

Go to **Components > Plug N Meet > Artifacts**:

1. Select a room and click **Show**.
2. Click an artifact ID to open its details page: room analytics table, per-user table (join/leave times, talk time, connection quality), and meeting summary where available.
3. Use **Download** for file artifacts, **Download Excel** for analytics reports, and **Delete** to remove.

##### 6. Advanced Configuration

###### Design Customization

- **Global:** **System > Global Configuration > plugNmeet > Design** (logo, colors, background, custom CSS) plus **Copyright** (footer branding). These are the defaults for every room.
- **Per-Room:** Room edit form **Design Customisation** tab — any value set here overrides global for that room only; empty fields inherit global.

![Joomla-globalcust.png](https://www.plugnmeet.org/img/Joomla/Joomla-globalcust.png)

###### Media Quality Defaults

Set global defaults under **System > Global Configuration > plugNmeet > Advanced**: `h720` webcam, `h1080fps15` screen share, `music` audio preset, and `vp8` codec suit most sites. Override per room in its **Advanced Settings** tab when needed. Raise only if your server and clients have bandwidth.

###### Client Load

Under **System > Global Configuration > plugNmeet**, **Client load** controls where the conference renders:

- **Remote (default):** embedded in Joomla (`view=room&layout=conference`). Branding comes from global + per-room Design. Use this unless you have a reason not to.
- **Redirect:** users are sent to `your-plugnmeet-server/?access_token=...`. The Joomla layout, global/room Design, and `?returned=true` handling no longer apply — branding is the server default.

###### Copyright Footer

Under **System > Global Configuration > plugNmeet > Copyright**, enable **Display copyright text** and set **Copyright text** (HTML allowed) to brand the client footer globally. This applies to every room created while enabled.

---

##### 7. Security and Privacy

- **Data Control:** You control the server (self-hosted or private cloud), so you control your data.
- **On-Premises Hosting:** The open-source plugNmeet server can run on your own hardware for GDPR/HIPAA-sensitive workloads.
- **Encryption:** All media is encrypted via WebRTC standards. Enable **E2EE** per room for maximum security.
- **Least Privilege:** Keep **Join as admin** and **Recording delete** restricted to Manager/Administrator/Super Users.

---

##### Common Questions & Troubleshooting

**Why do I see an 'Invalid API' error on my meeting page?**
You are still using the default demo credentials. Replace them under **System > Global Configuration > plugNmeet** with your own server details.

**Join fails with "same password" error on save?**
Moderator and attendee passwords must differ. Change one and save again.

**"Meeting has not started yet" for attendees?**
The room has **Moderator join first** enabled. A moderator must join first with the moderator password.

**How do I find the global settings after installation?**
As an administrator, go to **System > Global Configuration > plugNmeet**.

**Frontend shows no recordings?**
The user needs **Recording view** permission for that room (Permissions tab), and at least one cloud recording must exist for the room ID.

<a id="page-43"></a>

## اتصال آموزشی با LTI

LTI راه اتصال سامانه‌های مدیریت یادگیری به اتاق را فراهم می‌کند. تنظیمات URL و اعتبارنامهٔ ابزار در LMS ثبت می‌شود و پارامترهای سفارشی قابلیت‌های جلسه را تغییر می‌دهند. صفحه مثال‌های Canvas، Chamilo و Moodle را دارد. نقش ارسالی LMS باید به نقش اتاق نگاشت شود و مقادیر دقیق تنظیمات از نسخهٔ پیوست همین بخش برداشته شوند.

این بخش صفحهٔ `user-guide/lti` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/user-guide/lti) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/user-guide/lti.md) · [نسخهٔ محلی کامل](docs/upstream/docs/user-guide/lti.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Integrating with any LMS using LTI

Learning Tools Interoperability (LTI) is a standard protocol that allows you to securely connect Plug-N-Meet to virtually any Learning Management System (LMS) like Canvas, Moodle, Chamilo, and more.

This method is ideal if there is no dedicated plugin for your platform or if your institution restricts the installation of custom plugins. It allows instructors and students to join meetings and view recordings directly from their course page.

---

##### Part 1: General LTI Configuration Details

This section contains the information needed for **any** LTI integration.

###### Required Credentials

To connect any LTI-compatible platform, you will need the following three pieces of information. The **Consumer Key** and **Shared Secret** must be obtained from your own PlugNmeet server's API settings.

```
Launch URL: https://your-plug-n-meet.com/lti/v1
Consumer key: [Your plugNmeet API Key]
Shared secret: [Your plugNmeet API Secret]
```

For temporary testing, you can use the public demo server:

```
Launch URL: https://demo.plugnmeet.com/lti/v1
Consumer key: plugnmeet
Shared secret: EXAMPLE_API_SECRET
```

**Note:** It is highly recommended to allow the LTI tool to send the `launcher's name`; otherwise, the user will appear in the meeting with an empty name.

###### User Interface Overview

- **Instructor/Moderator View:** When a user with an instructor role launches the LTI tool, they will see an administrative view allowing them to manage sessions and view recordings.

![lti1.png](https://www.plugnmeet.org/img/lti/lti1.png)

- **Student/Attendee View:** When a student launches the tool, they will see a simple page to join the current live session.

![lti3.png](https://www.plugnmeet.org/img/lti/lti3.png)

###### Optional Custom Parameters

You can pass these optional parameters during the LTI configuration in your LMS to customize the meeting room experience.

| Field                 | Type    | Description                                                                 |
| :-------------------- | :------ | :-------------------------------------------------------------------------- |
| `room_duration`         | Number  | Session duration in minutes. `0` means no limit.                            |
| `allow_polls`           | Boolean | `true` (default) or `false`                                                 |
| `allow_shared_note_pad` | Boolean | `true` (default) or `false`                                                 |
| `allow_breakout_room`   | Boolean | `true` (default) or `false`                                                 |
| `allow_recording`       | Boolean | `true` (default) or `false`                                                 |
| `allow_rtmp`            | Boolean | `true` (default) or `false`                                                 |
| `mute_on_start`         | Boolean | `false` (default) or `true`                                                 |
| `primary_color`         | String  | Interface primary color in hex code (e.g., `#004D90`).                      |
| `secondary_color`       | String  | Interface secondary color in hex code.                                      |
| `background_color`      | String  | Interface background color in hex code.                                     |
| `custom_logo`           | String  | A direct HTTPS link to a logo image (e.g., `https://mydomain.com/logo.png`). |

---

##### Part 2: LMS Integration Examples

This section provides step-by-step examples for popular Learning Management Systems.

###### Canvas LMS Example

This example demonstrates how to add Plug-N-Meet as an external LTI app in Canvas.

1.  As an administrator, navigate to **Settings > Apps** in your Canvas course.
2.  Click **View App Configurations** and then **+App**.
3.  For **Configuration Type**, select **Manual Entry**.
4.  Fill in the fields with your PlugNmeet LTI credentials (Name, Consumer Key, Shared Secret, and Launch URL).
5.  Click **Submit**. The PlugNmeet app will now be available as an external tool for assignments and modules.

![lti-canvas-1.png](https://www.plugnmeet.org/img/lti/lti-canvas-1.png)

###### Chamilo LMS Example

1.  As an administrator, go to the **Plugins** page and enable the **IMS/LTI** plugin.
2.  Navigate to a course and click the pencil icon to edit its settings.
3.  Go to the **IMS/LTI** section and click **Configure external tools**.
4.  Click **Add external tool** and fill in the form with your PlugNmeet LTI credentials.
5.  After saving, **plugNmeet** will appear as an available tool in the **Interaction** section of your course.

![lti-chamilo-8.png](https://www.plugnmeet.org/img/lti/lti-chamilo-8.png)

###### Moodle LMS Example

This LTI method is a great alternative if your Moodle hosting provider does not allow the installation of custom activity plugins.

1.  As an administrator, go to **Site administration > Plugins > Manage tools**.
2.  Click **Configure a tool manually**.
3.  Fill in the form with your PlugNmeet LTI credentials (Tool name, Tool URL, Consumer key, Shared secret).
4.  Click **Save changes**.
5.  Now, any course creator can add this to their course. Turn on **Edit mode** in a course, click **Add an activity or resource**, and select **External tool**.
6.  From the **Preconfigured tool** dropdown, select the **plugNmeet** tool you just created.
7.  Give the activity a name, save, and it will be ready for your students.

![lti-moodle-3.png](https://www.plugnmeet.org/img/lti/lti-moodle-3.png)
![lti-moodle-4.png](https://www.plugnmeet.org/img/lti/lti-moodle-4.png)
![lti-moodle-5.png](https://www.plugnmeet.org/img/lti/lti-moodle-5.png)
![lti-moodle-12.png](https://www.plugnmeet.org/img/lti/lti-moodle-12.png)
![lti-moodle-13.png](https://www.plugnmeet.org/img/lti/lti-moodle-13.png)
---

##### Common Questions & Troubleshooting

**What is the difference between using LTI and the dedicated Moodle plugin?**
The dedicated Moodle plugin offers deeper integration with Moodle features like grading and permissions. However, the LTI method is a universal solution that works for nearly any LMS and is a great alternative if you cannot install custom plugins.

**Why do I get an 'Authentication Failed' error?**
This almost always means the **Consumer Key** or **Shared Secret** is incorrect. Double-check that you have copied these values correctly from your PlugNmeet server's API settings and that they are entered correctly in your LMS.

<a id="page-44"></a>

## راهنمای مدیر جلسه

مدیر می‌تواند میکروفن و دوربین، اشتراک صفحه، ضبط ابری و محلی، وایت‌برد و یادداشت مشترک را مدیریت کند. اتاق انتظار، قفل قابلیت‌ها، بی‌صداکردن کاربران، پیام خصوصی، نظرسنجی و اتاق فرعی به سیاست اتاق وابسته‌اند. پایان جلسه برای همه با خروج یک مدیر از اتاق دو عمل جدا هستند. ضبط ابری به ضبط‌کنندهٔ سرور و آماده‌شدن فایل بعد از پردازش وابسته است. قابلیت‌های AI، SIP و ورود جریان زنده باید پیش از نمایش به کاربر در سرور فعال باشند.

این بخش صفحهٔ `user-guide/moderator` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/user-guide/moderator) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/user-guide/moderator.md) · [نسخهٔ محلی کامل](docs/upstream/docs/user-guide/moderator.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Moderator Guide: How to Host and Manage a Meeting

Welcome to the Moderator Guide. This document provides a complete overview of all the features and controls available to you as a meeting host or moderator.

The moderator has full control over the meeting, including managing participants, using presentation tools, and configuring room settings. This guide is structured to help you quickly find the feature you need.

---

##### First Steps: Entering the Room & Choosing Your Audio

When you first enter a PlugNmeet room, you will be prompted to choose how you want to join the audio session. This is your most important first step.

![Onboarding-Desktop.jpg](https://www.plugnmeet.org/img/moderator/Onboarding-Desktop.jpg)

-   **Enable Microphone and Camera:** Choose this option to speak and share your video in the meeting. Your browser will ask for permission to access your microphone and camera.
-   **Continue as a listener:** Choose this option to join without a microphone or camera. You will be able to hear everything but will not be able to speak or share video.

**Important:** This choice is not permanent. If you decide to **Continue as a listener**, you can easily activate your microphone and camera later from the main control bar.

---
 
##### Section 1: Core Media Controls and recording

These are the basic controls for your audio and video, located in the bottom-left of the control bar.

###### Managing Your Microphone & Audio

- **Mute/Unmute:** If you joined with microphone and camera access, click the **Microphone** icon in the bottom-left of the footer to mute or unmute yourself.
- **Activate Microphone:** If you joined in "Continue as a listener" mode, the icon will show a headset. Click it to leave the listen-only session and activate your microphone.
- **Change Audio Source:** Click the small arrow next to the microphone icon to switch between different microphones.

![microphone-pop.webp](https://www.plugnmeet.org/img/moderator/microphone-pop.webp)

###### Managing Your Webcam & Virtual Backgrounds

- **Start/Stop Webcam:** Click the **Webcam** icon in the bottom-left of the footer to start or stop sharing your video. You can also toggle the visibility of your video in the meeting.
- **Virtual Backgrounds:** When starting your webcam, a preview window will appear. Here you can select a virtual background (blur or an image) or upload your own image (using the plus icon) before clicking **Share**.

![share-webcam-pop.webp](https://www.plugnmeet.org/img/moderator/share-webcam-pop.webp)

###### Recording the Meeting

Plug-N-Meet offers two distinct methods for recording your session: **Cloud Recording** (server-side) and **Local Recording** (browser-based). To start either, click the **REC** icon in the main control bar and choose your preferred option.

![recording.png](https://www.plugnmeet.org/img/moderator/recording.png)

###### Cloud Recording (Recommended)

This is the most reliable method. The recording is processed on your server, capturing the entire session without depending on your local computer's performance.

-   **Requirement:** To use Cloud Recording, the free and open-source `plugNmeet-recorder` component must have been installed on your server. This is an option during the main installation process.
-   **Troubleshooting:** If you see a "No recorder available" message, it means this component is not running on your server.
-   **How it works:** Simply select "Cloud Recording" to start. The REC icon will pulse to indicate that the server is recording. Click it again to stop.

    **New Recording Options:**
    When initiating a cloud recording, you now have additional options to optimize the recorded layout:
    *   **Enable auto-close chat panel**: If enabled, the public chat panel will automatically hide after a period of inactivity. This provides more screen space for webcams, screen shares, or presentations in the final recording.
    *   **Duration after last message (in minutes)**: Specify the duration of inactivity (in minutes) after which the chat panel will automatically close.

    These options ensure that your recordings are clean and focused on the primary content, especially useful for RTMP broadcasting where screen real estate is crucial. Our recording philosophy ensures that what you see is what gets recorded, so controlling the chat panel visibility directly impacts the final video.

###### Local Recording (Browser-Based)

This method uses your web browser's capabilities to capture and save the meeting directly to your computer. It has very specific requirements and is **only supported in Google Chrome**.

> **Important:** Local Recording is an advanced feature. For the best results, please read the following instructions carefully.

To capture all audio within the session (not just your own microphone), the recording must be performed by a "virtual" second user.

**Step-by-Step Instructions:**

1.  **Open a new Chrome window** (not just a new tab).
2.  Using the new window, **log in to the same meeting with a different user account** that has moderator privileges. You will now have two instances of yourself in the meeting.
3.  In this second window, click the **REC** icon and choose **Local Recording**.
4.  A Chrome prompt will appear asking you to share a tab. Select the tab of your *original* meeting session.
5.  **Crucial Step:** At the bottom of the Chrome prompt, you **must** check the box that says **"Share tab audio"**. If you do not enable this, only the video will be recorded, with no sound.

**Alternative Method (for capturing your microphone only):**

If you only need to record what your microphone picks up along with the main video, you can start local recording without a second user. However, you must **share your microphone first**, and *then* start the local recording. If you start recording before sharing your mic, your own audio will be excluded.

---

##### Section 2: Presentation & Collaboration Tools

These tools allow you to present content and collaborate with your attendees.

###### Sharing Your Screen

- Click the **Share Screen** icon in the main footer control bar.
- Your browser will prompt you to choose whether to share your entire screen, a specific application window, or a browser tab.
- To stop sharing, click the **Share Screen** icon again.

![share-screen.webp](https://www.plugnmeet.org/img/moderator/share-screen.webp)

###### Using the Interactive Whiteboard

The whiteboard is a powerful tool for real-time collaboration, powered by Excalidraw. To open it, click the **Whiteboard** icon in the main footer control bar.

![white-board.webp](https://www.plugnmeet.org/img/moderator/white-board.png)

###### Drawing & Annotation Tools
The toolbar provides a rich set of tools for drawing and brainstorming:
- Use the **pen tool** to draw freehand.
- Add **shapes**, **arrows**, and **text** to create diagrams.
- Change **colors**, **stroke styles**, and more.
- You can **export** the entire whiteboard canvas as a PNG or SVG file at any time.

###### Adding Content to the Whiteboard
There are two distinct methods for adding content to the whiteboard: one for images and one for office files.

**1. Adding Images**
You can add images (like JPG or PNG) to the whiteboard in several easy ways:
- **Drag and Drop:** Simply drag an image file from your computer and drop it directly onto the whiteboard canvas.
- **Copy and Paste:** Copy an image (from a website or a file on your computer) and paste it onto the whiteboard using `Ctrl+V` (or `Cmd+V` on Mac).
- **Toolbar Icon:** Click the **Image** icon in the whiteboard's toolbar to open a file selection window.

![upload-images.webp](https://www.plugnmeet.org/img/moderator/upload-images.webp)

**2. Adding Office Files & PDFs**
To present documents like PDFs, Word files, or PowerPoint presentations, you must use the dedicated file uploader.
- Click the **Files** icon in the whiteboard's toolbar.
- A new panel will open. From here, you can upload your office documents or PDF files.
- Once uploaded, the file will be listed below the upload option. From there, you can select the file you want to display on the whiteboard.

![upload-files.webp](https://www.plugnmeet.org/img/moderator/upload-files.webp)

**Note**
You cannot upload images using the "Files" uploader. It is specifically for office documents and PDFs.


###### Using the Shared Notepad

- **Activate/Hide:** Open the **More Options** menu (...) in the main footer control bar and select **Enable shared notepad**.
- The shared notepad allows all participants to write and edit text together in real-time. You can format text and import/export the content.

![shared-notepad.webp](https://www.plugnmeet.org/img/moderator/shared-notepad.webp)

###### Using AI Assistant

- **AI Activation:** Open the **More Options** menu (...) in the main footer control bar and select **AI Tools**.
- **AI Settings:** The **Insights AI settings** popup will appear. Here you can enable or disable the AI assistant, and choose to make it available for all participants or specific users.

![ai-tools.webp](https://www.plugnmeet.org/img/moderator/ai-tools.webp)

- **Show/Hide:** Click the **AI Chat** icon in the main footer control bar to show or hide it.
- **Send Message:** You can send a message to the AI assistant and get an answer.

![ai-chat-widget.webp](https://www.plugnmeet.org/img/moderator/ai-chat-widget.webp)

- **AI Meeting summarization:** Use this feature to summarize the meeting conversation, identify all key decisions, and create a list of action items.

![ai-meeting-summerization.webp](https://www.plugnmeet.org/img/moderator/ai-meeting-summerization.webp)

###### Transcription & Translation
- **Activation:** Open the **More Options** menu (...) in the main footer control bar and select **Transcription & Translation**.
- **Settings:** The Transcription and Translation Settings popup will appear. Here you can enable or disable Transcription & Translation.

![transcriptions-pop.webp](https://www.plugnmeet.org/img/moderator/transcriptions-pop.webp)

- **Start Transcription & Translation:** Click the **Transcription & Translation** Icon in the main footer control bar to start the transcription & translation.

![transcriptions-settings-pop.webp](https://www.plugnmeet.org/img/moderator/transcriptions-settings-pop.webp)

- **History:** You can view the **Live Transcription & Translation** in real-time. Click the **History** icon to see the full transcript history.

![transcriptions-history.webp](https://www.plugnmeet.org/img/moderator/transcriptions-history.webp)

**Important:** You must **Enable Microphone** to use **Transcription & Translation**.

###### Sharing an External Video

You can share a video from a URL (like YouTube) or a local file.

- Open the **More Options** menu (...) in the main footer control bar and select **Play External Video**.
- Paste a video URL or upload a video file (MP4, WebM).
- The video will appear for all participants, and you will have playback controls.

![external-video-player.webp](https://www.plugnmeet.org/img/moderator/external-video-player.webp)

###### Displaying an External Website (iframe)

You can display any website that allows embedding within an iframe.

- Open the **More Options** menu (...) in the footer control bar and select **Share External Link**.
- Enter the URL of the website you wish to share.

![external-link-iframe.webp](https://www.plugnmeet.org/img/moderator/external-link-iframe.webp)

---

##### Section 3: Managing Participants

These features allow you to control and interact with your attendees.

###### Viewing the Participant List

- Click the **Participants** icon in the footer control bar to open the participant panel. Here you can see everyone in the meeting.

###### Using the Waiting Room

If the waiting room is enabled, you must approve attendees before they can join.

- Open the **More Options** menu (...) in the footer control bar and select **Manage Waiting Room**.
- In the waiting room panel, you can see a list of users waiting to enter. You can **Accept** or **Reject** them individually, or all at once.
- You can also approve users directly from the participant list.

![manage-waiting-room.webp](https://www.plugnmeet.org/img/moderator/manage-waiting-room.webp)

###### Sending Private Chat Messages

- Open the **Participants** panel.
- Click the three-dot menu next to a user's name and select **Send private message**.
- This will open a private chat tab in the chat panel.

![private-message.webp](https://www.plugnmeet.org/img/moderator/private-message.webp)

###### Using Room Lock Settings

This feature allows you to restrict what attendees can do.

- Open the **More Options** menu (...) in the footer control bar and select **Meeting Lock Settings**.
- From here, you can lock attendees' microphones, webcams, screen sharing, whiteboard sharing, notepad sharing, and chat abilities.
- These locks do not affect other moderators.

![lock-settings.webp](https://www.plugnmeet.org/img/moderator/lock-settings.webp)

###### Muting All Users

- To mute all attendees at once, open the **More Options** menu (...) and select **Mute all users**.
- This will not mute other moderators.

---

##### Section 4: Advanced Engagement Features

Use these tools to create a more interactive and engaging session.

###### Creating and Managing Polls

- Open the **More Options** menu (...) in the footer control bar and select **Enable Polls**. This will enable polls for the meeting and you will see a **Polls** icon in the main footer control bar.
- Click **Create new Poll**.
- Define your question and answer options, then click **Create poll**.
- Attendees will be prompted to vote. You can view the results in real-time, end the poll, and publish the results to the public chat.

![polls.webp](https://www.plugnmeet.org/img/moderator/polls.webp)

###### Using Breakout Rooms

Split participants into smaller, temporary groups for focused discussions.

- Open the **More Options** menu (...) in the footer control bar and select **Manage Breakout Rooms**.
- Configure the number of rooms, the duration, and assign participants (or use random assignment).
- Click **Start breakout rooms**.
- While rooms are active, you can broadcast messages to all groups, join specific rooms, or end the breakout session.

![breakout-room.webp](https://www.plugnmeet.org/img/moderator/breakout-room.webp)

###### Live Streaming via RTMP

You can broadcast your meeting live to platforms like YouTube or Facebook.

- Open the **More Options** menu (...) in the footer control bar and select **Start Live Stream**.
- Enter the **Stream Key** and **Stream URL** provided by your streaming platform.
- Click **Start Streaming** to begin streaming.

![live-streaming.webp](https://www.plugnmeet.org/img/moderator/live-streaming.webp)

###### Live stream input
###### Stream Like a Pro: How to Bring OBS into Your plugNmeet Room
Getting started is incredibly easy. Here’s how you do it from within your plugNmeet room:

- Click the three-dot menu in the **top-right** corner and select **Settings**.
- Find **Live Stream Input**: In the settings panel, navigate to the **Live Stream Input** tab.
- Choose **Your "Channel" Type (RTMP or WHIP)**: You'll see an option called "Input Type." This lets you choose the technology for your stream. Think of it like choosing between two types of cables to connect your camera. RTMP is for professional streaming, while WHIP is for live streaming.
- Give **Your Stream a Name**: In the **Display Name** field, you can type in a name for your stream. This is the name that will appear in the participant list. For example, "Main Presentation" or "Guest Speaker Cam." If you leave it blank, it will default to "Broadcaster."
- Generate **Your Private Link**: Click the **Generate link** button. plugNmeet will instantly create a unique and secure address for your private stream. You'll see two fields: a URL and a Secret (or "Stream Key").

![live-stream-input.webp](https://www.plugnmeet.org/img/moderator/live-stream-input.webp)
![stream-key-host.webp](https://www.plugnmeet.org/img/moderator/stream-key-host.webp)

###### Connecting Your Streaming Software (like OBS)

Now, you just need to tell your software (we'll use OBS as an example) where to send the stream.

###### For RTMP (The Universal Standard):

- In OBS settings, go to **Stream**.
- For **Service**, choose **Custom...**
- Copy the **URL** from plugNmeet and paste it into the **Server** field in OBS.
- Copy the **Secret** from plugNmeet and paste it into the **Stream Key** field in OBS.

![obs-custom.webp](https://www.plugnmeet.org/img/moderator/obs-custom.webp)

###### For WHIP (The Modern Choice):

- In OBS settings, go to **Stream**.
- For **Service**, choose **WHIP**.
- Combine the **URL** and **Secret** from plugNmeet into a single line. For example: `https://your-plugnmeet-domain.com/whip/A_LONG_SECRET_KEY_HERE`
- Paste this full, combined address into the **Server** field in OBS.
- Leave the **Bearer Token** field in OBS blank.

![obs-whip.webp](https://www.plugnmeet.org/img/moderator/obs-whip.webp)

Once you hit **"Start Streaming"** in OBS, your stream will magically appear as a new participant in your plugNmeet room for everyone to see! You can checkout our blog post for more details. [Stream Like a Pro: How to Bring OBS into Your plugNmeet Room](https://www.plugnmeet.org/blog/obs-rtmp-whip-ingress)

###### SIP/VoIP Dial-In

**Server Configuration Required**
This is an advanced feature that will only be available if your system administrator has configured the SIP integration on the server.


This feature allows participants to join the meeting's audio by dialing a standard telephone number.

![sip_1.png](https://www.plugnmeet.org/img/moderator/sip_1.png)

- Click the three-dot menu in the **top-right** corner and select **Settings**.
- Navigate to the **SIP/VoIP Dial-In** tab.
- **Enable Service:** If the service is not already running, click the **Activate** button to start the SIP gateway for this session.

![sip_2.png](https://www.plugnmeet.org/img/moderator/sip_2.png)

- **View Dial-In Information:** Once enabled, this panel will display the **Dial-in Number(s)** and the unique **PIN** for the current meeting.
- **Share in Chat:** Click the **Share to Chat** button to automatically post the dial-in number and PIN to the public chat, making it easy for all participants to see.

---

##### Section 5: Session Management

These actions control the overall meeting session.

###### Application Settings

- Click the three-dot menu in the **top-right** corner and select **Settings**.
- Here you can change the application language, manage data-saving preferences, toggle notification sounds, and manage Live Stream settings.

![app-settings.webp](https://www.plugnmeet.org/img/moderator/app-settings.webp)

###### Ending vs. Logging Out of a Meeting

- **Leave Meeting:** If you select **Leave Meeting** from the top-right menu, you will leave the meeting, but the session will remain active for other participants.
- **End Meeting:** If you select **End meeting**, the session will be terminated for **everyone**.

![meeting-end.webp](https://www.plugnmeet.org/img/moderator/meeting-end.webp)

<a id="page-45"></a>

## اتصال Moodle

افزونه به شکل activity در دوره نصب می‌شود. تنظیمات اتصال در مدیریت Moodle و مجوز نقش‌ها در دوره تعیین می‌شوند. گزارش حضور تجمعی و تکمیل فعالیت از analytics استفاده می‌کند. ابزارهای بارگذاری وایت‌برد، نظرسنجی از بانک سؤال، اتاق فرعی از گروه‌های درس و پایش جلسهٔ فعال نیز در صفحه توضیح داده شده‌اند. ضبط‌های چندبخشی و artifact بعد از جلسه جداگانه مدیریت می‌شوند.

این بخش صفحهٔ `user-guide/moodle-integration` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/user-guide/moodle-integration) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/user-guide/moodle-integration.md) · [نسخهٔ محلی کامل](docs/upstream/docs/user-guide/moodle-integration.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### How to Add a Video Conferencing Activity to Moodle with Plug-N-Meet

Welcome to the definitive guide for integrating **plugNmeet**, a powerful, open-source **Moodle video conferencing** platform, into your Moodle courses. This document provides a comprehensive walkthrough, from initial setup to leveraging advanced features, ensuring you can build a seamless and engaging **virtual classroom** experience. With this **free video conference software**, you can host **live online classes**, webinars, and meetings without ever leaving your Moodle site.

---

##### 1. Introduction: Why plugNmeet is the Best Choice for Moodle

plugNmeet is not just another video meeting tool; it's a complete **virtual classroom** solution built for deep integration with Moodle.

-   **Seamless Native Experience**: The classroom loads directly within your Moodle layout. To your students and teachers, it feels like a natural part of your site, not a third-party tool. This boosts user adoption and reduces confusion.
-   **Open Source and Free**: As a leading **open source video conferencing** solution, plugNmeet gives you complete ownership. You control your data, your infrastructure, and your costs.
-   **Deep Moodle Integration**: plugNmeet hooks directly into Moodle's core systems, including the calendar, activity completion, grading, and user permissions, for a truly integrated experience.
-   **Rich Feature Set**: It's packed with tools designed for education, including an interactive whiteboard, breakout rooms, live polling, and advanced analytics.
-   **Unmatched Customization**: Brand the entire virtual classroom—from logos and colors to the layout itself—to perfectly match your institution's identity.

---

##### 2. Installation and Configuration

To get started with plugNmeet, your Moodle administrator must install the plugin and connect it to a plugNmeet server. This section covers the installation and configuration of the Moodle plugin itself.

###### Step 2.1: Install the plugNmeet Moodle Plugin

1.  Download the latest version of the plugin from the [Moodle Plugin Directory](https://moodle.org/plugins/mod_plugnmeet) or the [official GitHub repository](https://github.com/mynaparrot/moodle-mod_plugnmeet).
2.  Log in to your Moodle site as an administrator.
3.  Navigate to `Site administration > Plugins > Install plugins`.
4.  Upload the plugin's ZIP file and click **"Install plugin from the ZIP file"**.
5.  Moodle will check the plugin and prompt you to **"Upgrade Moodle database now"**. Proceed with the upgrade.

###### Step 2.2: Configure the Plugin Connection

This is the most critical step. You need to tell Moodle how to connect to your plugNmeet server. Your server administrator will provide you with the necessary credentials.

1.  After installation, you will be taken to the settings page. You can also find it later at `Site administration > Plugins > plugNmeet`.
2.  **API Configuration**:
    -   **Server URL**: Enter the full URL of your plugNmeet server (e.g., `https://pnm.yourdomain.com`).
    -   **API Key**: Enter the `PLUGNMEET_API_KEY` provided by your server administrator.
    -   **API Secret**: Enter the `PLUGNMEET_SECRET` provided by your server administrator.
3.  **Client Loading Method**:
    -   Choose **Remote (Embedded)**. This is the recommended setting for a seamless, white-label experience where the classroom loads inside your Moodle page.
    -   `Redirect` is a fallback that sends users to the plugNmeet server's default interface, which is not ideal for branding.
4.  **Save changes**. The plugin will attempt to connect to the server. If successful, you're ready to go! If not, double-check your credentials and server URL.

![Screenshot of the Moodle admin settings page with fields for API credentials and the client loading method highlighted.](https://www.plugnmeet.org/img/moodle/moodle_1.png)

---

##### 3. Creating and Managing a Live Session

Adding a video conference to a course is simple.

1.  Navigate to the Moodle course where you want to add the session.
2.  Turn **Edit mode** on.
3.  In the desired section, click **Add an activity or resource**.
4.  Select **plugNmeet** from the activity chooser. You will be taken to the activity settings page.

![add plugnmeet activity](https://www.plugnmeet.org/img/moodle/moodle_2.png)

###### Detailed Configuration Options

This form gives you granular control over every aspect of the session.

-   **General Settings**:
    -   **Room Title**: The name of your session (e.g., "Weekly Biology Seminar"). This is what students will see.
    -   **Welcome Message**: A custom message that appears in the chat when users join. Great for posting an agenda or instructions.
    -   **Max Participants**: Limit the number of users who can join the session.

-   **Room Features (Collapsible Section)**:
    -   **Moderator must join first**: Prevents students from starting the session without a teacher present.
    -   **Mute on start**: A great classroom management tool to ensure a quiet start to the session.
    -   **Allow Screen Share / Webcams**: Enable or disable these core features. You can even set it to "Admin only webcams" for a presentation-style webinar.
    -   **Allow Breakout Rooms / Max Breakout Rooms**: Lets moderators split the session into smaller groups. Set the maximum number of rooms you may need.
    -   **Pre-assign breakout rooms from course groups**: When enabled, each Moodle course group becomes a breakout room automatically — the room title is the group name and members are pre-assigned. Requires plugNmeet server v2.5.0+ and Moodle plugin v3.0.13+. Takes effect from the next session start; the moderator can still rearrange rooms live.
    -   **Allow Polls**: Enables live polls and quizzes in the session. Required if you want to push polls from Moodle quiz questions during a live class (see below).

-   **Recording Features**:
    -   **Allow Recording**: Master switch to enable/disable all recording types.
    -   **Allow Cloud Recording**: Enables recording directly to the server. This is the most common and useful option.
    -   **Enable Auto Start Cloud Recording**: The session will automatically start recording as soon as the first person joins.

-   **Default Lock Settings**:
    -   These settings establish the *initial* state of the room locks for participants (not moderators).
    -   For a secure-by-default classroom, you can **Lock Screen Sharing**, **Lock Whiteboard**, and **Lock Shared Notepad**. Participants will have to ask for permission (raise their hand) to use these tools.

-   **Availability**:
    -   **Available from**: Sets the start time for the session. This automatically creates an event in the Moodle Calendar, making it a crucial tool for scheduling your **Moodle webinar** or class.
    -   **Deadline**: Sets the end time. After this time, the "Join" button will disappear.

---

##### 4. In-Room Features for Engaging Online Classes

plugNmeet is packed with tools to foster interaction and collaboration.

-   **Interactive Whiteboard**: A powerful tool for real-time collaboration. You can draw, write text, upload and annotate PDFs/presentations, and work together with students.
-   **Screen Sharing**: Share your entire screen, a single application window, or a browser tab with high-quality resolution.
-   **Breakout Rooms**: Split your main session into smaller, private groups for focused discussions or project work. Moderators can move between rooms and broadcast messages to all. With Moodle course groups, rooms can be pre-assigned automatically (see below).
-   **Polls**: Create polls on the fly to engage students, check for understanding, or vote on topics. Results can be published live. Teachers can also push a poll straight from a Moodle quiz or question bank question into the live room (see below).
-   **Shared Notepad**: A simple, real-time collaborative text editor for taking shared notes.
-   **AI-Powered Features**: If enabled on the server, you can access:
    -   **Live Transcription & Translation**: Get a real-time transcript of the conversation, which can be translated by users into their preferred language.
    -   **AI Meeting Summarization**: Automatically generate and download a summary of the meeting's key points.

---

##### 5. Advanced Analytics, Attendance, and Completion Tracking

plugNmeet's integration with Moodle provides unparalleled insights into student engagement.

###### Cumulative Attendance Reporting

Unlike tools that only report on a single session, plugNmeet aggregates data across *all* sessions within a single Moodle activity.

-   **For Teachers**: Access the **Attendance** report from the activity page to see a summary of each student's total participation. This includes total time attended, whether they raised their hand, sent chat messages, etc.
-   **For Students**: Students have their own self-service dashboard where they can see their attendance status and track their own engagement against the activity's requirements.

![Screenshot of the teacher's attendance report, showing a list of students with their status and participation metrics.](https://www.plugnmeet.org/img/moodle/moodle_3.png)
![Screenshot of the students's attendance report, showing their participation metrics.](https://www.plugnmeet.org/img/moodle/moodle_4.png)

###### Advanced Activity Completion

This is one of the most powerful features for **Moodle virtual classrooms**. You can set highly specific, automated completion rules based on active participation.

1.  In the activity settings, go to the **Activity completion** section.
2.  Set **Completion tracking** to "Show activity as complete when conditions are met".
3.  Check the boxes for the criteria you want to enforce. Examples:
    -   Require students to attend for a cumulative total of **45 minutes**.
    -   Require students to **send at least one chat message**.
    -   Require students to **enable their webcam** at least once.

This system moves beyond passive viewing and encourages active, measurable engagement.

![Screenshot of the activity completion settings, with custom plugNmeet rules like "completionminutes" and "completionchatmessages" checked.](https://www.plugnmeet.org/img/moodle/moodle_5.png)

---

##### 6. Real-Time Session Monitoring and Interaction

###### In-Activity Live Dashboard

Once a session is live, teachers and moderators gain access to a real-time dashboard directly on the Moodle activity page. This powerful tool provides an at-a-glance overview of the session's status and allows for direct interaction. You can instantly see:

*   **Live Participant List**: A list of all users currently in the session.
*   **Activity Status**: See who has their webcam or microphone active.
*   **Session Summary**: A summary of key metrics, including participant count, recording status, and active media.

This live data allows moderators to manage the classroom effectively without having to constantly switch between windows.

![Screenshot of the Active live sessions with participant details.](https://www.plugnmeet.org/img/moodle/moodle_7.png)

###### Upload to Whiteboard from Moodle

While a session is running, administrators and teachers can upload files directly to the whiteboard from Moodle's file repository. This feature streamlines the process of sharing course materials without leaving the virtual classroom.

1.  From the live dashboard on the Moodle activity page, click the new **"Upload to Whiteboard"** button.
2.  This will open the familiar Moodle file picker.
3.  You can select an existing file from your course files.
4.  Once selected, the file is automatically added to the whiteboard's file list inside the live session, ready for you to present to the class.

![Screenshot of the Moodle file picker being used to upload a file to the whiteboard.](https://www.plugnmeet.org/img/moodle/moodle_11.png)

###### Add Poll from Quiz / Question Bank (Live)

While a session is running, teachers can push a poll into the live room straight from a Moodle quiz or question bank question — powered by the [Create Poll API](https://www.plugnmeet.org/docs/api/room/create-poll). Requires plugNmeet server v2.5.0+ and Moodle plugin v3.0.13+.

1.  From the live dashboard on the Moodle activity page, click **"Add Poll from Quiz"** (visible only while the room is active).
2.  Choose the source: a course **Quiz** or the **Question Bank** category.
3.  Pick a supported question (single-answer multiple choice or true/false).
4.  Choose the mode: **Quiz** (correct answer revealed after the poll ends) or **Vote** (responses only), plus **Anonymous voting** and **Auto-close after** (minutes, max 60).
5.  Submit — the poll appears in the live room immediately and participants can vote right away. The moderator manages it like any other poll (close, publish results, reopen).

Poll voting also feeds Moodle completion: enable **Voted in Poll** in the activity completion settings to require participation.

###### Pre-Assigned Breakout Rooms from Course Groups

Instead of sorting students manually, let Moodle groups do the work. Requires plugNmeet server v2.5.0+ and Moodle plugin v3.0.13+.

1.  In the activity settings, under **Room Features**, enable **Pre-assign breakout rooms from course groups**.
2.  Make sure the course has groups with members (e.g. `Site administration` or course `Participants > Groups`).
3.  From the next session start, each group becomes a breakout room: the room title is the group name and members are automatically assigned. Users with no group stay unassigned for the moderator to place live.
4.  The moderator can still rearrange, remove, broadcast, extend, or end rooms from the in-meeting breakout management panel.

Note: room details are fixed at session start, so group changes apply to the next session, not the running one.

###### Site-Wide Active Session Report (Admin Data Pull)

For administrators, plugNmeet offers a real-time "live data pull" of all active sessions across the entire Moodle site. This is an essential tool for monitoring site-wide activity.

-   **How to Access**: An administrator can add a custom link to their admin dashboard or simply navigate to `https://youmoodlesite.com/mod/plugnmeet/index.php`.
-   **What it Shows**: The report displays a list of all currently running rooms, showing the room title, course, number of participants, and session duration. This allows admins to see at a glance how the **Moodle video conferencing** system is being used.

![Screenshot of the Active Rooms Report, showing a table of live sessions with participant counts.](https://www.plugnmeet.org/img/moodle/moodle_6.png)

---

##### 7. Managing Post-Session Recordings and Artifacts

All files generated during your sessions are neatly organized and accessible directly from the Moodle activity page.

###### Recordings Tab
Here you can find all cloud recordings. Teachers can:
-   Play the recording.
-   Control visibility (show/hide from students).
-   Download the recording file.
-   Delete the recording.

![Screenshot of the Active Recordings.](https://www.plugnmeet.org/img/moodle/moodle_8.png)

###### Merge Multiple Recordings

If a session was recorded in multiple segments, you can now easily merge them into a single, continuous video directly from the Moodle interface.

1.  Navigate to the **Recordings** tab for the plugNmeet activity.
2.  You will see checkboxes next to each recording. Select the recordings in the specific order you want them to appear in the final video.
3.  Click the **"Merge Recordings"** button that appears at the bottom of the list.
4.  A confirmation pop-up will show the list of recording IDs you are about to merge, verifying the sequence of your selection.
5.  Confirm the action to start the merge process. The new, merged recording will appear in the list once it's ready.

![Screenshot of the Moodle recordings list with checkboxes and a merge button, showing the merge confirmation pop-up.](https://www.plugnmeet.org/img/moodle/moodle_12.png)

###### Artifacts Tab
This tab contains all other session-generated files, governed by Moodle permissions. This includes:
-   AI-generated meeting summaries.
-   Transcription files (e.g., VTT, SRT).
-   Detailed per-session analytics reports in JSON format.

![Screenshot of the Artifacts.](https://www.plugnmeet.org/img/moodle/moodle_9.png)

![Screenshot of the Analytics.](https://www.plugnmeet.org/img/moodle/moodle_10.png)

---

##### 8. Security and Privacy

plugNmeet is designed with security in mind for educational institutions.

-   **Role-Based Permissions**: The plugin respects Moodle's roles. For example, only users with the `mod/plugnmeet:viewrecordings` capability can see the recordings tab. The Active Rooms Report is also protected by a specific capability (`mod/plugnmeet:viewactiveroomsreport`).
-   **Secure by Default**: Use the **Default Lock Settings** in your activity to create a controlled environment where students must request permission to share their screen or whiteboard.
-   **End-to-End Encryption (E2EE)**: For maximum security, you can enable E2EE for your sessions.
-   **Guest Access Control**: Guest access can be disabled globally from the admin settings, and each session's guest link has a configurable expiration time.

This concludes the detailed guide to the **plugNmeet Moodle plugin**. You are now ready to build a secure, engaging, and fully-branded **Moodle video conferencing** experience for your users.

<a id="page-46"></a>

## مسیر راهنمای کاربران

این بخش مسیر ورود به راهنمای مدیر جلسه، شرکت‌کننده و اتصال‌های آمادهٔ CMS/LMS را مشخص می‌کند. برای آی‌نشست، راهنمای مدیر و شرکت‌کننده مبنای آموزش کاربران است و راهنماهای افزونه‌ها الگوی تنظیم نقش، دسترسی و خروجی پس از جلسه را توضیح می‌دهند.

این بخش صفحهٔ `user-guide/overview` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/user-guide/overview) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/user-guide/overview.md) · [نسخهٔ محلی کامل](docs/upstream/docs/user-guide/overview.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### Plug-N-Meet User Guide

Welcome to the Plug-N-Meet User Guide. This section provides comprehensive guides for everyone using the platform.

Whether you are a **moderator** managing a live session, an **attendee** joining a meeting, or an **administrator** integrating Plug-N-Meet with a platform like WordPress, Moodle, or Joomla, you will find the information you need here.

---

##### In-Meeting Experience

Learn how to use the features inside a Plug-N-Meet room. The capabilities depend on your role.

- **[Moderator Guide](https://www.plugnmeet.org/docs/user-guide/moderator):** A complete guide for hosts on how to manage meetings, participants, recordings, and use advanced features like breakout rooms and polling.
- **[Attendee Guide](https://www.plugnmeet.org/docs/user-guide/attendee):** A simple guide for participants on how to join meetings and use features like chat, screen sharing, and raising your hand.

---

##### Platform Integration Guides

Step-by-step instructions for integrating Plug-N-Meet with popular platforms.

- **[WordPress Integration](https://www.plugnmeet.org/docs/user-guide/wordPress-integration):** How to add video conferencing to your WordPress site.
- **[Moodle Integration](https://www.plugnmeet.org/docs/user-guide/moodle-integration):** How to add a video conference activity to your Moodle courses.
- **[Joomla Integration](https://www.plugnmeet.org/docs/user-guide/joomla-integration):** How to add video meetings to your Joomla site.
- **[Integrating with any LMS using LTI](https://www.plugnmeet.org/docs/user-guide/lti):** A guide for connecting to Canvas, Chamilo, or any other LTI-compatible platform.

<a id="page-47"></a>

## اتصال WordPress

افزونه با URL سرور، API key و secret متصل می‌شود. مدیر اتاق می‌سازد، نقش و مجوز حضور و تنظیمات ضبط و ظاهر را تعیین می‌کند و اتاق را با shortcode یا لینک مستقیم منتشر می‌کند. پس از جلسه، ضبط و artifact از بخش مدیریت قابل دریافت‌اند. صفحه همچنین روش بارگذاری کلاینت، جایگزینی قالب و hookهای توسعه را توضیح می‌دهد. این مسیر برای نصب WordPress است و نمونهٔ استفاده از API در یک محصول میزبان را نشان می‌دهد.

این بخش صفحهٔ `user-guide/wordPress-integration` را پوشش می‌دهد. [صفحهٔ رسمی](https://www.plugnmeet.org/docs/user-guide/wordPress-integration) · [سورس ثابت](https://github.com/mynaparrot/plugNmeet-docs/blob/44a474a17d965e17ebb68c43991205768ef4c27f/docs/user-guide/wordPress-integration.md) · [نسخهٔ محلی کامل](docs/upstream/docs/user-guide/wordPress-integration.md)

### متن فنی، جدول‌ها و مثال‌های منبع

#### How to Add Video Conferencing to WordPress with the Plug-N-Meet Plugin

Welcome to the definitive guide for integrating **plugNmeet**, the powerful, open-source video conferencing solution, into your WordPress website. This document provides a complete walkthrough, from installation to leveraging advanced features, enabling you to host secure, feature-rich video meetings directly from your own domain. With plugNmeet, you can create a seamless virtual meeting experience that feels like a native part of your site.

---

##### 1. Key Features of the WordPress Integration

The plugNmeet plugin is more than just a way to embed a video call; it's a deep integration that brings powerful management tools directly into your WordPress dashboard.

*   **Seamless Shortcode Integration:** Embed a full-featured video conference room into any page or post with a simple WordPress shortcode.
*   **Backend Management:** Create rooms, manage recordings, view analytics, and access meeting artifacts without ever leaving your WordPress admin area.
*   **Advanced Post-Session Analytics:** Go beyond simple attendance. Get detailed reports on participant engagement, talk time, and connection quality, all exportable to Excel.
*   **Full Asset Control:** Easily access, download, and manage all room assets, including recordings, speech transcripts, and AI-generated meeting summaries.
*   **Role-Based Permissions:** Use standard WordPress user roles to control who can join meetings as a moderator and who can manage recordings.
*   **Developer Friendly:** Extend the plugin's functionality by hooking into webhook events to process real-time meeting data.

---

##### 2. Installation & Initial Setup

This section covers installing the plugin and connecting it to your PlugNmeet server.

###### Step 2.1: Install and Activate the Plugin

1.  From your WordPress dashboard, navigate to **Plugins > Add New**.
2.  In the search bar, type **plugnmeet**.
3.  Find the official Plug-N-Meet plugin and click **Install Now**, then **Activate**.


###### Step 2.2: Connect to Your Server (API Settings)

This is the most important step. The plugin requires a connection to a working PlugNmeet server.

1.  From your WordPress dashboard, navigate to the new **Plug-N-Meet** menu item and select **Settings**.
2.  Fill in the three fields: **Server URL**, **API Key**, and **API Secret** with the credentials from your self-hosted or PlugNmeet Cloud server.
3.  **You must replace the default demo credentials for the plugin to work.**

**Note:** The plugin includes pre-configured demo credentials to help you test its features immediately. This demo server is a shared resource and is **not intended for production use** as it can be unreliable. For any important meetings, we strongly recommend using your own server to ensure a stable and professional experience for you and your users.

###### Client Load Method

In the settings, you will also find the **Client Load** option. This important setting changes *how* the video meeting is displayed to your users and how you can customize its appearance.

*   **Remote (Recommended for most users):** This option loads the video conferencing interface directly from your main plugNmeet server and embeds it seamlessly into your WordPress page. It is the simplest method and always uses the most up-to-date software.
    *   **Customization:** You can easily customize the branding (logos, colors, etc.) for each room using the **Design Customization** tab when creating or editing a room.

*   **Redirect:** This option sends the user to a separate browser tab to join the meeting directly on your plugNmeet server. It is simple and keeps the meeting interface separate from your website's layout.
    *   **Customization:** Branding is controlled by the **default configuration on your plugNmeet server**. Any changes made in the "Design Customization" tab within WordPress will **not** apply with this method.

*   **Local (Advanced, for developers):** This option is for users who need to make deep, structural changes to the meeting interface. It requires you to download the client application files, modify the source code, and host them on your own server. This provides maximum control but requires web development experience and manual updates.

![API settings page for the PlugNmeet WordPress video conferencing plugin.](https://www.plugnmeet.org/img/wordpress/1.png)

**IMPORTANT: Refresh Your Permalinks**
After saving your settings for the first time, you **must** refresh your WordPress permalinks. This ensures that the direct links to your meeting rooms (e.g., `yoursite.com/plugnmeet/room/your-room-id`) will work correctly and not result in a "404 Not Found" error.

1.  Go to **Settings > Permalinks** in your WordPress dashboard.
2.  You do not need to change any settings. Simply click the **"Save Changes"** button.

You must repeat this step anytime you change the **Room Slug Path** in the Plug-N-Meet settings.


---

##### 3. Creating & Publishing Your First Meeting

###### Step 3.1: Navigating to the Room Creation Page

To create a new meeting room, navigate to **Plug-N-Meet > Rooms** from your WordPress dashboard and click the **"Add New"** button. This will take you to the detailed room configuration page, which is organized into several tabs.

![Navigating to the 'Add New' room page in the PlugNmeet WordPress conferencing plugin.](https://www.plugnmeet.org/img/wordpress/2.png)

###### Step 3.2: Detailed Room Configuration

This form gives you granular control over every aspect of the session.

![The detailed room configuration interface for the PlugNmeet WordPress plugin, showing multiple tabs for customization.](https://www.plugnmeet.org/img/wordpress/3.png)

###### Basic
This tab covers the essential information for your room.
- **Title & Description:** Give your meeting room a clear title and an optional description.
- **Passwords:** Set separate passwords for moderators and attendees.
- **Welcome Message:** Customize the initial message that appears in the chat when users join.
- **Max Participants:** Set a limit on the number of users who can join (use `0` for unlimited).

###### Room Features
This section allows you to enable or disable the core interactive features of the meeting room, such as webcams, screen sharing, and chat.

###### Other Features
This tab provides a wide range of settings to control specific functionalities within the room:
- **Recording:** Enable or disable cloud recording.
- **External Broadcasting:** Configure options for streaming your meeting to external platforms.
- **Chat, Shared Notepad, Whiteboard:** Toggle these collaborative tools on or off.
- **External Media Player:** Allow moderators to play media from external sources.
- **Waiting Room:** Enable a waiting room to screen participants before they join.
- **Breakout Rooms:** Allow moderators to split participants into smaller groups.
- **Display External Link:** Configure a custom link to be displayed to participants.
- **Ingress:** Set up options for bringing external media streams into the meeting.
- **Polls:** Enable the polling feature for interactive feedback.
- **SIP Dial-in:** Configure options for users to join via a SIP phone call.
- **End-to-End Encryption:** Enable E2EE for maximum security.

###### Insights AI
In this tab, you can enable or disable AI-powered features for the meeting, such as automatic transcription and meeting summaries.

###### Default Lock Settings
Configure the initial "locked" state for various features when a participant joins. For example, you can start the meeting with microphones, webcams, or the whiteboard locked to maintain order. Moderators can unlock these features during the meeting.

###### Design Customization
This tab allows you to override the global design settings for this specific room. You can change the logo, background, and color scheme to match your brand or the topic of the meeting.

###### Permission
This powerful tab allows you to define exactly how different WordPress user roles interact with the meeting room and its recordings.

![The user role permission table in the PlugNmeet WordPress video conferencing plugin, allowing granular control over meeting access.](https://www.plugnmeet.org/img/wordpress/4.png)

**Smart Defaults:** The plugin starts with sensible defaults. Roles with administrative capabilities (like `edit_users` or `edit_posts`, which typically apply to Administrators and Editors) are automatically set to join as **moderators** with full access to recordings. All other roles, including non-logged-in **Guests**, default to joining as **attendees** with more restricted permissions.

**Customization:** These defaults are just a starting point. You can override them for each role using the following options:
- **Join as Moderator/Attendee:** Determines the user's in-room privileges.
- **Require Password:** If checked, users with this role must enter the attendee password to join, even if they are logged in.
- **Allow View/Play/Download/Delete Recordings:** Provides granular control over how each role can interact with the room's recordings. The "Guest/Public" role cannot be granted permission to delete recordings.

###### Step 3.3: Publish the Room

You have two primary methods for making a room available to your users.

###### Option A: Embed in a Page with a Shortcode

This method is ideal for integrating a meeting room into an existing page or post alongside other content.

1.  Go to **Plug-N-Meet > Rooms** and find the room you just created.
2.  In the **Shortcode** column, copy the unique shortcode for that room.
3.  Paste the shortcode into the content area of any WordPress page or post.
4.  Publish the page. When a user visits, they will see a customizable form to join the meeting.

![Copying a room shortcode to embed a self-hosted video conference in a WordPress page with the PlugNmeet plugin.](https://www.plugnmeet.org/img/wordpress/5.png)

###### Option B: Use a Clean, Shareable Direct Link

Every room you create automatically has a direct, shareable link. This is perfect for sending in emails or for a clean, professional look. The URL structure consists of a base slug and the unique room ID: `yoursite.com/[slug]/[room-id]`.

You have full control over this structure:

1.  **Customize the URL Slug (Optional):**
    -   Go to **Plug-N-Meet > Settings > Options**.
    -   Find the **Room Slug Path** field. By default, it is `plugnmeet/room`.
    -   You can change this to something more unique, like `meetings`. Your new URL structure will be `yoursite.com/meetings/your-room-id/`.

2.  **Ensure a Clean Layout with a Host Page:**
    -   By default, visiting the direct link may show the meeting inside your theme's standard post layout (with sidebars, etc.). To ensure a clean, distraction-free view, you should assign a blank page as a template.
    -   **Create a Blank Page:** Go to **Pages > Add New** and create a new, empty page. Name it "Meeting Host" or something similar.
    -   **Assign the Host Page:** Go back to **Plug-N-Meet > Settings > Options** and select the page you just created from the **Room Host Page** dropdown.

3.  **Refresh Permalinks:**
    -   After customizing the slug or setting the host page, you **must** refresh your permalinks.
    -   Go to **Settings > Permalinks** and simply click **"Save Changes."**

Now, all your direct room links will open in a clean, professional layout using your custom URL structure.

---

##### 4. Post-Session Management

After a meeting concludes, all generated assets are neatly organized and accessible from the **Plug-N-Meet > Rooms** page in your WordPress admin area. From the rooms list, you can access the **Recordings** and **Artifacts** for each session.

###### Managing Recordings

![The recordings management page in the PlugNmeet WordPress plugin, showing options to download and merge video files.](https://www.plugnmeet.org/img/wordpress/7.png)

The **Recordings** page for each room allows you to:
- Download and delete recordings.
- **Merge multiple recording segments** into a single, continuous MP4 file for easy distribution.

###### Accessing Artifacts

![The artifacts management page in the PlugNmeet WordPress plugin, showing downloadable analytics reports and AI summaries.](https://www.plugnmeet.org/img/wordpress/6.png)

The **Artifacts** page contains all other files and data generated during the session. This includes:
- **Analytics Reports:** Download comprehensive post-meeting analytics in Excel format, detailing participant engagement, talk time, connection quality, and more.
- **AI-Generated Summaries:** Access and download meeting summaries created by the AI.
- **Transcription Files:** Download speech-to-text transcripts of the meeting audio.

![Artifact details](https://www.plugnmeet.org/img/wordpress/8.png)

---

##### 5. Advanced Configuration

###### Overriding Plugin Templates for Theme Developers

For maximum flexibility, the plugNmeet plugin allows theme developers to override any of its front-end template files by creating copies in their own theme directory. The plugin follows the standard WordPress template hierarchy, searching for files in this order:

1.  `[your-child-theme]/plugnmeet/`
2.  `[your-parent-theme]/plugnmeet/`
3.  `[plugin-directory]/public/partials/` (the default)

**How to Override:**

1.  **Locate the original file:** Find the template you want to change inside the plugin's `/wp-content/plugins/plugnmeet/public/partials/` directory.
2.  **Replicate the structure:** In your theme's folder (preferably a child theme), create a `plugnmeet` directory. If the original file is in a subfolder (e.g., `parts/`), create that subfolder inside your `plugnmeet` directory as well.
3.  **Copy and modify:** Copy the original file to the new location in your theme and customize it to your needs.

**Example 1: Minor Customization**
To change the layout of the join form, you would copy:
`[plugin-dir]/public/partials/parts/login-form.php`
to
`[your-theme]/plugnmeet/parts/login-form.php`
and then edit the file in your theme.

**Example 2: Major Structural Changes**
To completely change the overall structure of the room view, you can override the main display file. Copy:
`[plugin-dir]/public/partials/plugnmeet-public-display.php`
to
`[your-theme]/plugnmeet/plugnmeet-public-display.php`
This allows you to add custom wrapper divs, change the order of template parts, or even conditionally load different parts based on your own logic.

###### Developer Hooks

The plugin fires a `plugnmeet_webhook_data` action hook, allowing you to process real-time event data from your meetings. You can hook into this in your theme's `functions.php` or a custom plugin.

```php
add_action( 'plugnmeet_webhook_data', 'my_plugnmeet_webhook_handler', 10, 1 );

function my_plugnmeet_webhook_handler( $data ) {
  // Process the webhook data here
  // e.g., log to a file, send to a CRM, etc.
}
```

---

##### 6. Security and Privacy

Security is a core strength of the plugNmeet platform.
*   **Data Control:** Because you control the server (either self-hosted or through our private cloud), you control your data.
*   **On-Premises Hosting:** For maximum privacy, the open-source plugNmeet server can be installed on your own on-premises hardware, ensuring sensitive data never leaves your network. This is ideal for organizations with strict data compliance needs (e.g., GDPR, HIPAA).
*   **Encryption:** All communication is encrypted by default using WebRTC standards. For maximum security, you can enable **End-to-End Encryption (E2EE)** for your meetings.

## حقوق استفاده و منشأ

متن رسمی این فایل از مستندات منتشرشدهٔ MynaParrot تحت مجوز MIT گردآوری شده است. راهنمای فارسی و توضیحات اتصال آی‌نشست در کنار آن قرار گرفته‌اند. فایل اصلی مجوز و اعلام حق مؤلف باید همراه نسخه‌های بازنشرشده حفظ شود. [مجوز اصلی بایگانی‌شده](docs/upstream/LICENSE)

MIT License

Copyright (c) 2022 MynaParrot

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
