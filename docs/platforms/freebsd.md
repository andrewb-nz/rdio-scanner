![Rdio Scanner](../images/rdio-scanner.png?raw=true)\

\pagebreak{}

# What is it ?

[Rdio Scanner](https://github.com/chuot/rdio-scanner) is an open source software that ingest and distribute audio files generated by various software-defined radio recorders. Its interface tries to reproduce the user experience of a real police scanner, while adding its own touch.

You can listen to [Rdio Scanner](https://github.com/chuot/rdio-scanner) on any modern browsers using the integrated web app.

# Recorders compatibility

[Rdio Scanner](https://github.com/chuot/rdio-scanner) works with any radio recorder, as long as they can create audio files separated by conversations or transmissions.

Here is a list of recorders known to work with [Rdio Scanner](https://github.com/chuot/rdio-scanner):

| Recorder                                                       | API | Dirwatch |
| -------------------------------------------------------------- | --- | -------- |
| [Trunk Recorder](https://github.com/robotastic/trunk-recorder) | X   | X        |
| [RTLSDR-Airband](https://github.com/szpajder/RTLSDR-Airband)   |     | X        |
| [SDRTrunk](https://github.com/DSheirer/sdrtrunk)               |     | X        |
| [voxcall](https://github.com/aaknitt/voxcall)                  | X   |          |
| [ProScan](https://www.proscan.org/)                            |     | X        |
| [DSDPlus Fast Lane](https://https://www.dsdplus.com/)          |     | X        |

# Improve your experience on the go

You can enjoy your [Rdio Scanner](https://github.com/chuot/rdio-scanner) on the go on your mobile device with the native app.

[![Available on the App Store](../images/app-store-badge.png?raw=true)](https://apps.apple.com/us/app/rdio-scanner/id1563065667#?platform=iphone)
[![Get it on Google Play](../images/google-play-badge.png?raw=true)](https://play.google.com/store/apps/details?id=solutions.saubeo.rdioScanner)

\pagebreak{}

# Quick start

ALWAYS DOWNLOAD THE LATEST VERSION OF [RDIO SCANNER](https://github.com/chuot/rdio-scanner) FROM ITS OFFICIAL REPOSITORY AT **[HTTPS://GITHUB.COM/CHUOT/RDIO-SCANNER](https://github.com/chuot/rdio-scanner)**. PRECOMPILED VERSIONS ARE LOCATED UNDER THE [RELEASES TAB](https://github.com/chuot/rdio-scanner/releases).

1. Make sure you have the latest version of [Rdio Scanner](https://github.com/chuot/rdio-scanner).

2. Extract the contents of the archive somewhere on your computer. The configuration file and the database file will be created in this same folder.

        rdio@pc-freebsd:~ $ mkdir rdio-scanner
        rdio@pc-freebsd:~ $ cd rdio-scanner
        rdio@pc-freebsd:~/rdio-scanner $ unzip \
        > ~/rdio-scanner-freebsd-amd64-v7.0.0.zip 
        Archive:  ../rdio-scanner-freebsd-amd64-v7.0.0.zip
         extracting: rdio-scanner  
         extracting: rdio-scanner.pdf

3. Run the executable.

        rdio@pc-freebsd:~/rdio-scanner $ ./rdio-scanner
        
        Rdio Scanner v7.0.0
        ----------------------------------
        2022/12/19 08:16:36 server started
        2022/12/19 08:16:36 main interface at http://pc-freebsd:3000
        2022/12/19 08:16:36 admin interface at http://pc-freebsd:3000/admin

4. Access the administrative dashboard to finalize the configuration.

> _Note that the default password is **rdio-scanner**_

![admin dashboard login](../images/admin-login.png?raw=true)\

5. Configure your [Rdio Scanner](https://github.com/chuot/rdio-scanner) instance.

![admin dashboard todos](../images/admin-todos.png?raw=true)\

\pagebreak{}

# Advanced customizations

## Listening on an alternate port

Here we want our [Rdio Scanner](https://guthub.com/chuot/rdio-scanner) instance to listen for connections on the standard HTTP port.

> Because privileged ports requires root privileges, we need to run the command as _root_.

        rdio@pc-freebsd:~/rdio-scanner $ su root
        Password:
        root@pc-freebsd:/home/rdio/rdio-scanner # ./rdio-scanner -listen :80
        
        Rdio Scanner v7.0.0
        ----------------------------------
        2022/12/19 08:19:38 server started
        2022/12/19 08:19:38 main interface at http://pc-freebsd
        2022/12/19 08:19:38 admin interface at http://pc-freebsd/admin

## Listening on a SSL port

It is recommended to share your [Rdio Scanner](https://github.com/chuot/rdio-scanner) instance over the internet by listening to an SSL port.

You can use your own SSL certificates that match your domain name with `-ssl_cert_file` and `-ssl_key_file`. If your certificate comes with an intermediate CA certificate, simply add its contents to the standard certificate file.

        rdio@pc-freebsd:~/rdio-scanner $ su root
        Password:
        root@pc-freebsd:/home/rdio/rdio-scanner # ./rdio-scanner       \
        > -listen :80                                                  \
        > -ssl_cert_file mycert.crt                                    \
        > -ssl_key_file mykey.key                                     \
        > -ssl_listen :443
        
        Rdio Scanner v7.0.0
        ----------------------------------
        2022/12/19 10:34:29 server started
        2022/12/19 10:34:29 main interface at http://pc-freebsd
        2022/12/19 10:34:29 main interface at https://pc-freebsd
        2022/12/19 10:34:29 admin interface at https://pc-freebsd/admin

If you don't want to worry about SSL certificates, you can use the built-in Let's Encrypt auto-cert feature. This requires that you have both port 80 (HTTP) and port 443 (HTTPS) open to the world. Also, your domain name should point to your IP address where [Rdio Scanner](https://github.com/chuot/rdio-scanner/) is running. The advantage of this approach is that everything is done automatically, no certificate request, no certificate renewal.

        rdio@pc-freebsd:~/rdio-scanner $ su root
        Password:
        root@pc-freebsd:/home/rdio/rdio-scanner # ./rdio-scanner       \
        > -listen :80                                                  \
        > -ssl_auto_cert mydomain.com                                  \
        > -ssl_listen :443

## Save your advanced configuration to a config file

You don't want to have to type everytime a long list of arguments. No problem, you can save your advanced configuration to a file by adding the **-config_save** argument.

        rdio@pc-freebsd:~/rdio-scanner $ ./rdio-scanner       \
        > -listen :80                                         \
        > -ssl_auto_cert mydomain.com                         \
        > -ssl_listen :443                                    \
        > -config_save
        2022/12/19 10:37:00 rdio-scanner.ini file created

All of your parameters passed as arguments to [Rdio Scanner](https://github.com/chuot/rdio-scanner) have been saved to an INI file which has the same arguments/values list.

        db_file = rdio-scanner.db
        db_type = sqlite
        listen = :80
        ssl_auto_cert = mydomain.com
        ssl_listen = :443

Then simply run [Rdio Scanner](https://github.com/chuot/rdio-scanner) without any arguments.

        rdio@pc-freebsd:~/rdio-scanner $ su root
        Password:
        root@pc-freebsd:/home/rdio/rdio-scanner # ./rdio-scanner
        
        Rdio Scanner v7.0.0
        ----------------------------------
        2022/12/19 10:38:28 server started
        2022/12/19 10:38:29 main interface at http://pc-freebsd
        2022/12/19 10:38:29 main interface at https://pc-freebsd
        2022/12/19 10:38:29 admin interface at https://pc-freebsd/admin

## Install Rdio Scanner as a service

Everything looks ok now, but you do not want to manualy start [Rdio Scanner](https://github.com/chuot/rdio-scanner) after every reboot. No problem, you can install your server instance as a service. **Note that you will be required to proceed in a command prompt that is running as Administrator**.

        rdio@pc-freebsd:~/rdio-scanner $ su root
        Password:
        root@pc-freebsd:/home/rdio/rdio-scanner # ./rdio-scanner \
        > -service install
        root@pc-freebsd:/home/rdio/rdio-scanner # ./rdio-scanner \
        > -service start
        root@pc-freebsd:/home/rdio/rdio-scanner # ps ax | grep rdio-scanner
        1522  -  S       0:00.06 /usr/home/rdio/rdio-scanner/rdio-scanner -service run
        1524  0  S+      0:00.00 grep rdio-scanner

\pagebreak{}

## Get the full list of arguments

To get the whole list of arguments you can pass to [Rdio Scanner](https://github.com/chuot/rdio-scanner), simply pass the argument **-h**.

        rdio@pc-freebsd:~/rdio-scanner $ ./rdio-scanner -h
        Usage of ./rdio-scanner:
          -admin_password string
                change admin password
          -base_dir string
                base directory where all data will be written
          -cmd string
                advanced administrative tasks (use -cmd help for usage)
          -config string
                server config file (default "rdio-scanner.ini")
          -config_save
                save configuration to rdio-scanner.ini
          -db_file string
                sqlite database file (default "rdio-scanner.db")
          -db_host string
                database host ip or hostname (default "localhost")
          -db_name string
                database name
          -db_pass string
                database password
          -db_port uint
                database host port (default 3306)
          -db_type string
                database type, one of sqlite, mariadb, mysql, postgresql (default "sqlite")
          -db_user string
                database user name
          -listen string
                listening address (default ":3000")
          -service string
                service command, one of start, stop, restart, install, uninstall
          -ssl_auto_cert string
                domain name for Let's Encrypt automatic certificate
          -ssl_cert_file string
                ssl PEM formated certificate
          -ssl_create
                create self-signed certificates
          -ssl_key_file string
                ssl PEM formated key
          -ssl_listen string
                listening address for ssl
          -version
                show application version
        rdio@pc-freebsd:~/rdio-scanner $ 

\pagebreak{}

# Need help ?

You can ask your questions or post your comments on the [Rdio Scanner Discussions](https://github.com/chuot/rdio-scanner/discussions) at **[https://github.com/chuot/rdio-scanner/discussions](https://github.com/chuot/rdio-scanner/discussions)**.

# Show your appreciation, support the author

If you like [Rdio Scanner](https://github.com/chuot/rdio-scanner), **[consider starring the GitHub repository](https://github.com/chuot/rdio-scanner/stargazers)** to show you appreciation to the author for his hard work. It cost nothing but is really appreciated.

If you use [Rdio Scanner](https://github.com/chuot/rdio-scanner) for commercial purposes or derive income from it, **[sponsor the project](https://github.com/sponsors/chuot)** to help support continued development.

[![Follow us on Twitter](../images/twitter-badge.png?raw=true)](https://twitter.com/RdioScanner)

**Happy Rdio scanning !**

\pagebreak{}