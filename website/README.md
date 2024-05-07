# A word on Svelte

When trying to build the webpage associated to the user service the question arised about which web framework to use. In the past, we used [react](https://react.dev/learn) for the [sogclient]([https://](https://github.com/Knoblauchpilze/sogclient)) project. This was quite powerful but also very complex to get to work.

After a bit of research, we stumbled upon the [Svelte](https://kit.svelte.dev/docs/introduction) framework: this seems to be relatively lightweight. That being said it seems to use some strange syntax to represent component, of course different than what react uses.

## Run the code

To start the website, the [installation instructions](https://kit.svelte.dev/docs/creating-a-project), to run the project we need to use `npm run dev`. The internet indicates that to open the corresponding browser tab we can use the more complete:

```bash
npm run dev -- --open
```

## Approach to the framework

Taken from the [documentation](https://kit.svelte.dev/docs/routing#page-page-svelte), it seems like the routing is done on the file system. A

Generally speaking, it seems like to build a page with a certain route, you can create:
* `path/to/page.ts`: this allows to fetch data from remote source before the page is loaded.
* `path/to/page.svelte`: this defines the html content of the page and can access whatever was returned from the `page.ts`.
* `path/to/layout.svelte`: allows to define a common layout applying to all pages below the current path and inheriting from the parent layouts.

We also have a bunch of files that we will want to explore to correctly set-up the actions (e.g. [.prettierrc](users-dashboard/.prettierrc)).
