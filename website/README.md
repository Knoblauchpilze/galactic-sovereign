# A word on Svelte

When trying to build the webpage associated to the user service the question arised about which web framework to use. In the past, we used [react](https://react.dev/learn) for the [sogclient](https://github.com/Knoblauchpilze/sogclient) project. This was quite powerful but also very complex to get to work.

After a bit of research, we stumbled upon the [Svelte](https://kit.svelte.dev/docs/introduction) framework: this seems to be relatively lightweight. That being said it seems to use some strange syntax to represent component, of course different than what react uses.

## Approach to the framework

Taken from the [documentation](https://kit.svelte.dev/docs/routing#page-page-svelte), it seems like the routing is done on the file system.

Generally speaking, it seems like to build a page with a certain route, you can create:

- `path/to/page.ts`: this allows to fetch data from remote source before the page is loaded.
- `path/to/page.svelte`: this defines the html content of the page and can access whatever was returned from the `page.ts`.
- `path/to/layout.svelte`: allows to define a common layout applying to all pages below the current path and inheriting from the parent layouts.

It is also a full-stack framework, meaning that it serves both the client side (i.e. the pages that the user will see on their browser) and the backend (i.e. the code to serve the pages that will be sent to users).

In out case the later part will be 'underutilized' as we anyway already have our backend in the form of the go application. That being said it is nice to save some work to server the pages.

# Run the code

## Development build

To start the website locally, the [installation instructions](https://kit.svelte.dev/docs/creating-a-project), to run the project we need to use `npm run dev`. The internet indicates that to open the corresponding browser tab we can use the more complete:

```bash
npm run dev -- --open
```

Alternatively a [Makefile](users-dashboard/Makefile) is provided with a `make dev` target.

This usecase is interesting when developing a new feature: svelte handles a local server and opens a tab in the browser to present the website as it will look like for users.

## Production build

Deploying the website so that it can be accessed by users requires two components:

- the server which will handle the requests and serve the pages.
- the pages that are served to the users.

It also requires quite a lot of other things as described in this [reddit post](https://www.reddit.com/r/webdev/s/gEqYH5T0pg) but this would be the basics.

This [svelte doc page](https://kit.svelte.dev/docs/adapter-node) explains how to package everything for a node environment. This [page](https://kit.svelte.dev/docs/adapters) provides more information about the process in general.

The idea is to first build the code into production-ready artifacts with:

```bash
npm run build
```

It is then possible to start the server locally with:

```bash
ORIGIN=http://localhost:3000 node /path/to/the/build/folder
```

The first part corresponds to the problem described in this [SO post](https://stackoverflow.com/questions/73790956/cross-site-post-form-submissions-are-forbidden). With this, the website can be accessed locally on the browser by going to `http://localhost:3000/dashboard/login`.

For convenience, the [Makefile](users-dashboard/Makefile) defines a `build` target to handle the build. Most likely though it is better to use the root [Makefile](/Makefile) and build the docker images for the webserver: this handles everything and allows to predictably setup the webserver.
