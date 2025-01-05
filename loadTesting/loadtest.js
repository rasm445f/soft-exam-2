import http from "k6/http";
import { sleep, check } from "k6";

export let options = {
  stages: [
    { duration: "1s", target: 500 }, // Ramp-up to ? users
    { duration: "4m", target: 500 }, // Hold ? users
    { duration: "1s", target: 0 }, // Ramp-down
  ],
};

export default function () {
  const res = http.get("http://localhost:8083/api/restaurants/1/menu-items"); // Change to your project's endpoint
  check(res, {
    "is status 200": (r) => r.status === 200,
    "response time < 200ms": (r) => r.timings.duration < 200,
  });
  sleep(1);
}
