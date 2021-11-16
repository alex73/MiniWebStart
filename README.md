MiniWebStart
============

This is replacement for Java Web Start since Oracle deprecated it.

MiniWebStart doesn't work as part of browser, but was created for simplify software installation and updates on the end-user desktops.

This application also improves some things that was not so good in the Java Web Start:

* It works not only with Java applications, but can handle any applications, even non-application files.

* It doesn't require to use only installed java on the client computer, but can be used for distribute java as part of your application. That means, you don't need to update java manually on the all client computers.

* It doesn't require jars signing - it's just an optional thing. Probably, certificate of your https server will be enough for validate application on the user side.

The main work principles are almost the same as in Java Web Start:
1. You need to compile your application for each client platform
2. Need to create some deployment description(like jnlp file)
3. Send MiniWebStart to the clients
4. Client will just run MiniWebStart, then MiniWebStart will download/update/start files as described in deployment descriptor

MiniWebStart is just a one executable file and don't need any installation. Client should just run it with some parameters, like deployment description URL.
Developer can build own version of MiniWebStart with custom built-in parameters for send to user. In this case, user will need just to click on the MiniWebStart file for start your application, and will not need to add any parameters.

TODO
====

MiniWebStart is fully working, but I'm planning to make some improvements:
1. GUI
2. File signing
3. Offline allowed execution

Deployment description
======================

Deployment description looks like jnlp file:

```
<?xml version="1.0" encoding="UTF-8"?>
<mws>
  [<information>]
    [<offline-allowed/>]
  [</information>]
  <resources [os="windows|linux|darwin"] [bits="32|64"] [minMemory="800m|2g"] [maxMemory="1500m|4g"]>
    <unpack href="jre-8.zip" [toDir="java/"] [useModes="true|false"] />
    <unpack ... />
    <file href="start.sh" [toFile="run.sh"] [mode="0755"]/>
    <file ... />
  </resources>
  <resources ...>
  </resources>
  <startup [os="windows|linux|darwin"] [bits="32|64"] [minMemory="800m|2g"] [maxMemory="1500m|4g"] file="./run.sh" />
  <startup ... />
</mws>
```

Optional attributes and tags are in square brackets. Minimum memory size is exclusive, maximum memory size - inclusive. M is Mebibytes, i.e. 1024\*1024 bytes, G is gibibytes, i.e. 1024\*1024\*1024 bytes.

Keep in mind that you need to use current directory prefix for execute downloaded files, like ".\start.cmd" or "./start.sh".

Execution parameters
====================

You can set command line parameter like "--remote=https://yousite/desc.xml", or define constant in the predefined.go for declare deployment description location. Other command line parameters will be sent to startup application.

Building
========

You can build application by "DOCKER_BUILDKIT=1 docker build --output=out ." command.
