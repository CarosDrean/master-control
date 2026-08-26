export const API_PORT = 8080;

export function getApiBase(): string {
    if (typeof window === "undefined") {
        return `http://localhost:${API_PORT}/master-control`;
    }
    return `http://${window.location.hostname}:${API_PORT}/master-control`;
}