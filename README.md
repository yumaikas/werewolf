# Werewolf: An Artisanal Web-based outliner

This is a Work in progress web outliner that I plan on using to help augment how I approach large systems and long term projects.

Planned features include being able to clone/mirror nodes across multiple places in the outline, making sure that the outliner works serviceably on mobile, and adding support for running Lua scripts to automate parts of working on the outline. 
I also want to use this as a tool for building and/or understanding software. I was inspired by a Python/Desktop based outliner (I forget the name) that features cloning outline nodes.  So, for now, I'm building it out slowly, and trying to get time in on it where it makes sense. 

Right now, this is pre-0.1, barely anything works other than displaying the outline, the tests are still in a messy state, and so on. Caveat Emptor. 
UPDATE: Now I can edit the contents of existing nodes. Next on the list is creating new nodes at arbitrary places in the outline.


## Crazy ideas for down the road

- Vim key bindings for editing text, and/or navigating nodes
- a TUI for editing the outline, so that the Web GUI isn't required
- Some sort of query language for getting lists of nodes

## Contributing

I'm not looking for other contributors at this point, but if this venture of re-inventing outliners interesting, send a PR with the list of things you'd like to work on, or raise an issue. 

## TODOs currently in the code (using this to manage what things I want to come back too)

Below is the output of running `rg TODO: -tgo -tjs` at the top level of the codebase. 
I'm going to try to use this as a deferred work tracker for the team of myself. 

I plan on taking a little time every day I sit down to work on this and try to knock off at least one or two TODOs, or make them clearer, so that people other than myself could stand a chance of understanding them.

```
app.go:// TODO: Also detect and update the meta attributes of the node, if present
app.go:	// TODO: Return node later
main.go:	// TODO: Create a real db init function
main.go:		// TODO: Build a custom 500 page recoverer at some point.
db.go:	// TODO: Test remove and reorder later
db.go:		// TODO: Try to fix this later, with more food in my body.
WIPOutlineScripts.js:// TODO: Each node should also have an expand/collapse button
WIPOutlineScripts.js:// TODO: Add UI for creating new nodes
WIPOutlineScripts.js:// TODO: Add UI for deleting nodes
WIPOutlineScripts.js:// TODO: Add UI for re-ordering and re-paretning nodes
WIPOutlineScripts.js:            // TODO: Clean this up on the server side
templates/templates.go:// TODO: Come back to this if I find it works better "gopkg.in/russross/blackfriday.v2"
templates/templates.go:// TODO: integrate "github.com/microcosm-cc/bluemonday" into this code if I start trusting more than one user.
```
