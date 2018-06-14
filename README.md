MiniWebStart
============

This is replacement for Java Web Start since Oracle deprecated it.

MiniWebStart doesn't work as part of browser, but was created for simplify software installation and updates on the end-user desktops.

This application also improves some things that was not so good in the Java Web Start:

* It works not only with Java applications, but can handle any applications, even non-application files.

* It doesn't require to use only installed java on the client computer, but can be used for distribute java as part of your application. That means, you don't need to update java manually on the all client computers.

* It doesn't require jars signing - it's just an optional thing. Probably, certificate of your https server will be enough for validate application on the user side.

The main work principles are almost the same as in Java Web Start:
1. You need to compile your application
2. Need to create some deployment description(like jnlp file)
3. Send MiniWebStart to the clinet.
4. Client will run MiniWebStart for download/update/start files as described in deployment descriptor

MiniWebStart is just a one executable file and don't need any installation. Client should just run it with some parameters, like deployment description URL.
Developer can build own version of MiniWebStart with custom built-in parameters for send to user. In this case, user will need just to click on the MiniWebStart file for start your application.

TODO
====
MiniWebStart is fully working, but I'm planning to make some improvements:
1. GUI is planed
2. Support jnlp-like file as deployment description
3. File signing
4. Extend jnlp syntax for support resources depends on OS bits flag and memory size value
