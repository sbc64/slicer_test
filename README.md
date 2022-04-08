Steps:

1. ~~Read raw tpc.
2. ~~open tcp on relay
3. ~~use readFrom
4. use dup2 to get a second tcp connection that the websocket will read.
5. Read url paramater projectId before opening tcp to relay. You can use raw tcp parsing to read the headers: https://github.com/gobwas/httphead
6. open webocket internally of the splicer and read the contents of websocket without sening payload responses.
