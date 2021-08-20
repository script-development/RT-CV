## Q and A

#### Whats a "rolling key"?

A key that is changed every time it is used and cannot be reversed.
This protects agains re-play attacks and stealing a key used as authentication is useless.

#### Whats the salt?

For a few reasons:
- To make it harder to obain information about the master key used
- For every salt you can create a new rolling key and thus allowing multiple services or threads to use the same master key to communicate with the server
- If something goes wrong for example one of the servers loses internet or an internal server error happends you can start over with a diffrent salt

#### Whats the seed?

The seed is created every time the server starts and protects agains replay attacks.
Every time the this app restarts the previouse seeds are gone and thus someone who recorded some elses earlier keys can re-use them and with this that's not possible