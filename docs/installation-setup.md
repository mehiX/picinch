# Installation and setup
## Before installation
1. Set up a host system with Docker and Docker-Compose installed. **Digital Ocean example**
2. Aquire a domain name, or add a sub-domain to a domain you already own. Set the `A` record for the domain or subdomain to the IP address of your server.
The following instructions assume a Ubuntu Server host. Other Linux distributions will be similar. **CentOS warning**
## Installation steps
These instructions assume a Ubuntu Server host with Docker installed. A basic installation requires the creation of three files on the server.
1. Create and edit `/srv/picinch/docker-compose.yml` as described below. This Docker Compose file specifies the PicInch and MariaDB containers to be downloaded from Docker Hub, the settings to run them on the host system, and essential application parameters.
2. Add  `/srv/systemd/service/picinch.service` This file defines PicInch  as a service in the host system, and sets it to run at system startup. Make it executable. I.e. `chmod 664 picinch.service`.
3. Run `systemctl start picinch.service` When issued the first time, this fetches PicInch and MariahDB containers from Docker Hub, and starts PicInch. Then PicInch sets up the database, creates the folders to hold images and certificates (in`/srv/picinch/`), On later invocations, PicInch starts immediately.
4. Connect to your server by domain name using a web browser and see that you can log in.
## After installation
5. Customise the appearance of the website using the optional files described below. Restart PicInch for changes to take effect. (See Commands.)
6. Log in to PicInch as administrator and add potential users with status set to `known`. Send invitations to the users, inviting them to sign-in.
7. If desired, make one or more `active` (signed-up) users an `admin` or `curator`.
8. Arrange for a regular backup of the database and photos. This may be an option from your hosting supplier, or you may need to do it regularly yourself.
9. Review the security of your system.