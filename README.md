# onefunc
This hideous grab bag of utility functions has no overall theme or vision, it is simply a collection of unrelated functions I find myself implementing repeatedly.

I cycled through multiple ways to handle this (writing functions locally, creating a github account with multiple projects, etc.) and settled on this as the least-ugly way to handle the issue of reducing replication of small, useful functions.

The initial idea was that every function would be isolated and self-contained ("onefunc") and you would not import the package, but instead the package would include a utility to copy the desired functions into your project. Ultimately I couldn't find any technical justification for handling it that way, it just seemed dumb.

The closest thing to a rule related to this project is: No dependencies. If a package needs a non-system dependency then it becomes a likely candidate to be spun off into a separate repo. _This_ repo should be an endpoint, not dragging in a huge amount of third-party code if you want to include it.

But don't include it! It __is__ pretty hideous, after all. If, by chance, you have the misfortune to stumble across this project and find something useful, my suggestion would be to copy it or start your own similar repo.
