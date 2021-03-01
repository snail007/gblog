# GBLOG

A blog engine based on [gmc](https://github.com/snail007/gmc) golang framework.

# RUNTIME REQUIREMENTS

## WITHOUT BLEVE
1. libc >=2.14 (debian 8+,ubuntu 14.10+,centos 7+)

## WITH BLEVE
1. libc >=2.18 (debian 8+,ubuntu 16.04+,centos 8+)
2. libstdc++ >=6.0.21 (debian 8+,ubuntu 16.04+,centos 8+)

# PREVIEW

[Demo](https://gblog-demo.herokuapp.com/)

[Demo Manage](https://gblog-demo.herokuapp.com/manage/) root 123456

[snail007's blog using gblog](https://www.host900.com/)

![](/doc/images/intro0.png)

![](/doc/images/intro1.png)

![](/doc/images/intro2.png)

![](/doc/images/intro3.png)

![](/doc/images/intro4.png)

![](/doc/images/intro5.png)

![](/doc/images/intro6.png)

# BUILD

Requirements
1. libc >=2.18 (debian 8+,ubuntu 16.04+,centos 8+)
2. libstdc++ >=6.0.21 (debian 8+,ubuntu 16.04+,centos 8+)
3. go>=1.16
4. docker installed

```shell script

find / | grep libc.so.6 
# such as, contains: 64, /lib/x86_64-linux-gnu/libc.so.6

strings /lib/x86_64-linux-gnu/libc.so.6 |grep GLIBC_2.18

find / | grep libstdc++.so.6
# such as, contains: 64, /usr/lib/x86_64-linux-gnu/libstdc++.so.6.0.22

ls -al /usr/lib/x86_64-linux-gnu/libstdc++.so.6
# lrwxrwxrwx 1 root root 19 Feb 15  2018 /usr/lib/x86_64-linux-gnu/libstdc++.so.6 -> libstdc++.so.6.0.21

```

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