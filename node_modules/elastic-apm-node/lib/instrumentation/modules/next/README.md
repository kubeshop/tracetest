# Next.js instrumentation dev notes

Instrumenting the Next.js servers (the dev and prod servers) involves shimming
a number of files. An overview of the instrumentation is provided here.
Some of the complexity here is that I found the Next.js server internals not
always ammenable to instrumentation. There isn't a clean separation between
"determine the route", "run the handler for that route", "expose an error
or the result".

There are a number of ways to deploy (https://nextjs.org/docs/deployment) a
Next.js app. This instrumentation works with "Self-Hosting", and using Next.js's
built-in server.

Here is the Next.js "server" class hierarchy:

    class Server (in base-server.ts)
        class NextNodeServer (in next-server.ts, used for `next start`)
            class DevServer (in dev/next-dev-server.js, used for `next dev`)


## dist/server/next.js

This is the first module imported for `require('next')`. It is used solely
to `agent.setFramework(...)`. Doing so in "next-server.js" can be too late
because it is lazily imported when creating the Next server -- by which point
metadata may have already been sent on the first APM agent intake request.


## dist/server/next-server.js

This file in the "next" package implements the `NextNodeServer`, the Next.js
"production" server used by `next start`. Most instrumentation is on this class.

The Next.js server is a vanilla Node.js `http.createServer` (http-only, https
termination isn't supported) using `NextNodeServer.handleRequest` as the request
handler, so every request to the server is a call to that method.

User routes are defined by files under "pages/". Generally, an incoming request path:
    GET /a-page
is resolved to one of those pages:
    ./pages/a-page.js
The user files under "./pages/" are loaded by `NextNodeServer.findPageComponents`.
**We instrument `findPageComponents` to capture the resolved page name to use
for the transaction name.**

There are also other built-in routes to handle redirects, rewrites, static-file
serving (e.g. `GET /favicon.ico -> ./public/favicon.ico`), and various internal
`/_next/...` routes used by the Next.js client code for bundle loading,
server-side generated page data, etc.  At server start, a call to
`.generateRoutes()` is called which returns a somewhat regular data structure
with routing data. **We instrument *most* of these routes to set
`transaction.name` appropriately for most of these internal routes.** A notable
limitation is the `public folder catchall` route that could not be cleanly
instrumented.

An error in rendering a page results in `renderErrorToResponse(err)` being
called to handle that error. **We instrument `renderErrorToResponse` to
`apm.captureError()` those errors.** (Limitation: There are some edge cases
where this method is not used to handle an exception. This instrumentation isn't
capturing those.)

*API* routes ("pages/api/...") are handled differently from other pages.
The `catchAllRoute` route handler calls `handleApiRequest`, which resolves
the URL path to a possibly dynamic route name (e.g. `/api/widgets/[id]`,
**we instrument `ensureApiPage` to get that resolve route name**), loads the
webpack-compiled user module for that route, and calls `apiResolver` in
"api-utils/node.ts" to execute. **We instrument that `apiResolve()` function
to capture any errors in the user's handler.**


## dist/server/dev/next-dev-server.js

This file defines the `DevServer` used by `next dev`. It subclasses
`NextNodeServer`. The instrumentation in this file is **very** similar to that
in "next-server.js". However, some of it needs to be repeated on the `DevServer`
class to capture results specific to the dev-server. For example
`DevServer.generateRoutes()` includes some additional routes.


## dist/server/api-utils/node.js

See the `apiResolve()` mention above.


