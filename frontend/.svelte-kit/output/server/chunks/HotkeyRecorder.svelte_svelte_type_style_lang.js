const API_PORT = 8080;
function getApiBase() {
  if (typeof window === "undefined") {
    return `http://localhost:${API_PORT}/master-control`;
  }
  return `http://${window.location.hostname}:${API_PORT}/master-control`;
}
const BASE_URL = getApiBase();
function appIconUrl(path) {
  return `${BASE_URL}/apps/icon?path=${encodeURIComponent(path)}`;
}
export {
  appIconUrl as a
};
