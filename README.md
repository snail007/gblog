# GBLOG

A blog engine based on [gmc](https://github.com/snail007/gmc) golang framework.

# PREVIEW

[Demo](https://gblog-demo.herokuapp.com/)

[Demo Manage](https://gblog-demo.herokuapp.com/manange/) root 123456

[snail007's blog using gblog](https://www.host900.com/)

![](/doc/images/intro0.png)

![](/doc/images/intro1.png)

![](/doc/images/intro2.png)

![](/doc/images/intro3.png)

![](/doc/images/intro4.png)

![](/doc/images/intro5.png)

![](/doc/images/intro6.png)

# BUILD
Building based on the docker, you must have docker installed.
Then run

```shell script
./pack.sh
```
After execute done, release directory location at : gblog-release.

Include 32bit & 64bit of windows, linux, mac.

# RUN

After you pack it, just run:

```shell
cd gblog-release/gblog-linux64-release
./gblog
```

visit http://`YOUR_IP`:6800/

# LOGIN

visit http://`YOUR_IP`:6800/manage/  
username: `root`  
password: `123456`  