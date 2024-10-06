export function GetHost(): string {
  // If origin is port 8080 then were on the dev server
  if (window.location.origin.includes("8080")) {
    return "http://localhost:3000";
  }

  // If origin is port 3000 then were on the prod server
  if (window.location.origin.includes("3000")) {
    return "http://localhost:3000";
  }

  // return origin
  return window.location.origin;
}
