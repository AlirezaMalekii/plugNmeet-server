import importlib.util
import unittest
from pathlib import Path


spec = importlib.util.spec_from_file_location("restore_haproxy", Path(__file__).with_name("restore-haproxy.py"))
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)


class SharedProxyRollbackTest(unittest.TestCase):
    def test_restores_only_owned_backends(self):
        previous = Path(__file__).with_name("haproxy.cfg").read_text()
        current = previous.replace("send-proxy-v2", "send-proxy")
        current = current.replace("frontend websites_http\n", "frontend websites_http\n    acl other_domain hdr(host) -i other.example.com\n")
        current += "\nbackend other_app\n    mode http\n    server other 127.0.0.1:18090\n"
        restored = module.restore(previous, current)
        self.assertIn("send-proxy-v2", restored)
        self.assertIn("acl other_domain hdr(host) -i other.example.com", restored)
        self.assertTrue(restored.endswith("backend other_app\n    mode http\n    server other 127.0.0.1:18090\n"))

    def test_missing_backend_rejects_without_output(self):
        previous = Path(__file__).with_name("haproxy.cfg").read_text()
        with self.assertRaises(ValueError):
            module.restore(previous, "backend turn\n    mode tcp\n")


if __name__ == "__main__":
    unittest.main()
