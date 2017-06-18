#!/usr/bin/env bash
#
#
#                  Corsair Installer Script
#
#   Quick hack by lindenstaub to work with systemd and make the folders/files expected by corsair
#
#   Original version:   https://github.com/shipy4rd/corsair/corsair-install.sh
#   Requires: bash, mv, rm, tr, type, curl/wget, tar (or unzip on OSX and Windows)
#
# This script safely installs Corsair into your PATH (which may require
# password authorization). Use it like this:
#
#	$ curl https://github.com/ship4rd/corsair | bash
#	 or
#	$ wget -qO- https://github.com/shipy4rd/corsair | bash
#
# If you want to get Corsair with extra plugins, use -s with a
# comma-separated list of plugin names, like this:
#
#	$ curl https://github.com/shipy4d/corsair | bash -s http.git,http.ratelimit,dns
#
# In automated environments, you may want to run as root.
# If using curl, we recommend using the -fsSL flags.
#
# !!! This script will probably break on any system without systemd. !!!
# Make a pull request if you have a patch to make it init system agnostic.
# https://github.com/lindenstaub/corsair-install.sh
#
#
#                      TROUBLESHOOTING
# Error: corsair.service start request repeated too quickly, refusing to start.
#
# Comment out this line in the /etc/systemctl/system/corsair.service file to see an actual error:
# `Restart=on-failure`
# Then you have to make systemctl see that the file is changed:
# `systemctl daemon-reload`
#
# Or to manually start corsair and see it's output directly
# `sudo -u www-data -h /usr/local/bin/corsair -log stdout -agree=true -conf=/etc/corsair/Corsairfile -root=/var/tmp`

install_corsair()
{
	trap 'echo -e "Aborted, error $? in command: $BASH_COMMAND"; trap ERR; return 1' ERR
	corsair_os="unsupported"
	corsair_arch="unknown"
	corsair_arm=""
	corsair_plugins="$1"
	install_path="/usr/local/bin"

	# Termux on Android has $PREFIX set which already ends with /usr
	if [[ -n "$ANDROID_ROOT" && -n "$PREFIX" ]]; then
		install_path="$PREFIX/bin"
	fi

	# Fall back to /usr/bin if necessary
	if [[ ! -d $install_path ]]; then
		install_path="/usr/bin"
	fi

	# Not every platform has or needs $sudo_cmd (see issue #40)
	((EUID)) && [[ -z "$ANDROID_ROOT" ]] && sudo_cmd="sudo"

	#########################
	# Which OS and version? #
	#########################

	corsair_bin="corsair"
	corsair_dl_ext=".tar.gz"

  # TODO would be nice to make an if to set which init script/system it's using
  # and later install it to the right place
  corsair_init_path="init/linux-systemd/"
  corsair_init="corsair.service"

	# NOTE: `uname -m` is more accurate and universal than `arch`
	# See https://en.wikipedia.org/wiki/Uname
	unamem="$(uname -m)"
	if [[ $unamem == *aarch64* ]]; then
		corsair_arch="arm64"
	elif [[ $unamem == *64* ]]; then
		corsair_arch="amd64"
	elif [[ $unamem == *86* ]]; then
		corsair_arch="386"
	elif [[ $unamem == *armv5* ]]; then
		corsair_arch="arm"
		corsair_arm="5"
	elif [[ $unamem == *armv6l* ]]; then
		corsair_arch="arm"
		corsair_arm="6"
	elif [[ $unamem == *armv7l* ]]; then
		corsair_arch="arm"
		corsair_arm="7"
	else
		echo "Aborted, unsupported or unknown architecture: $unamem"
		return 2
	fi

	unameu="$(tr '[:lower:]' '[:upper:]' <<<$(uname))"
	if [[ ${unameu} == *DARWIN* ]]; then
		corsair_os="darwin"
		corsair_dl_ext=".zip"
		vers=$(sw_vers)
		version=${vers##*ProductVersion:}
		IFS='.' read OSX_MAJOR OSX_MINOR _ <<<"$version"

		# Major
		if ((OSX_MAJOR < 10)); then
			echo "Aborted, unsupported OS X version (9-)"
			return 3
		fi
		if ((OSX_MAJOR > 10)); then
			echo "Aborted, unsupported OS X version (11+)"
			return 4
		fi

		# Minor
		if ((OSX_MINOR < 5)); then
			echo "Aborted, unsupported OS X version (10.5-)"
			return 5
		fi
	elif [[ ${unameu} == *LINUX* ]]; then
		corsair_os="linux"
	elif [[ ${unameu} == *FREEBSD* ]]; then
		corsair_os="freebsd"
	elif [[ ${unameu} == *OPENBSD* ]]; then
		corsair_os="openbsd"
	else
		echo "Aborted, unsupported or unknown os: $uname"
		return 6
	fi

	########################
	# Download and extract #
	########################

	echo "Downloading corsair for $corsair_os/$corsair_arch$corsair_arm..."
	corsair_file="corsair_${corsair_os}_$corsair_arch${corsair_arm}_custom$corsair_dl_ext"
	corsair_url="https://github.com/shipy4rd/corsair/$corsair_os/$corsair_arch$corsair_arm?plugins=$corsair_plugins"
	echo "$corsair_url"

	# Use $PREFIX for compatibility with Termux on Android
  echo $PREFIX
	rm -rf "$PREFIX/tmp/$corsair_file"

	if type -p curl >/dev/null 2>&1; then
		curl -fsSL "$corsair_url" -o "$PREFIX/tmp/$corsair_file"
	elif type -p wget >/dev/null 2>&1; then
		wget --quiet "$corsair_url" -O "$PREFIX/tmp/$corsair_file"
	else
		echo "Aborted, could not find curl or wget"
		return 7
	fi

	echo "Extracting bin..."
	case "$corsair_file" in
		*.zip)    unzip -o "$PREFIX/tmp/$corsair_file" "$corsair_bin" -d "$PREFIX/tmp/" ;;
		*.tar.gz) tar -xzf "$PREFIX/tmp/$corsair_file" -C "$PREFIX/tmp/" "$corsair_bin" ;;
	esac
	chmod +x "$PREFIX/tmp/$corsair_bin"

	echo "Extracting init script..."
	case "$corsair_file" in
		*.zip)    unzip -o "$PREFIX/tmp/$corsair_file" "$corsair_init_path$corsair_init" -d "$PREFIX/tmp/" ;;
		*.tar.gz) tar -xzf "$PREFIX/tmp/$corsair_file" -C "$PREFIX/tmp/" "$corsair_init_path$corsair_init" ;;
	esac
  echo $corsair_init_path$corsair_init


	# Back up existing corsair, if any
	corsair_cur_ver="$("$corsair_bin" --version 2>/dev/null | cut -d ' ' -f2)"
	if [[ $corsair_cur_ver ]]; then
		# corsair of some version is already installed
		corsair_path="$(type -p "$corsair_bin")"
		corsair_backup="${corsair_path}_$corsair_cur_ver"
		echo "Backing up $corsair_path to $corsair_backup"
		echo "(Password may be required.)"
		$sudo_cmd mv "$corsair_path" "$corsair_backup"
	fi

	echo "Putting corsair in $install_path (may require password)"
	$sudo_cmd mv "$PREFIX/tmp/$corsair_bin" "$install_path/$corsair_bin"
  # Took out if statement because it wasn't running even on systems where it should, not sure how this change breaks compatibility with other systems
	if setcap_cmd=$($sudo_cmd which setcap); then
    $sudo_cmd setcap cap_net_bind_service=+ep "$install_path/$corsair_bin"
	fi
	$sudo_cmd rm -- "$PREFIX/tmp/$corsair_file"

  # TODO make this init system sensitive
  init_script_path="/etc/systemd/system/"
  echo "Putting init script in $init_script_path"
  $sudo_cmd cp $PREFIX/tmp/$corsair_init_path$corsair_init $init_script_path$corsair_init
  $sudo_cmd rm -rf $PREFIX/tmp/$corsair_init_path
  $sudo_cmd chown root:root $init_script_path$corsair_init
  $sudo_cmd chmod 644 $init_script_path$corsair_init
  $sudo_cmd systemctl daemon-reload

  # check init script installation
  # TODO make this init system neutral
  # I couldn't make this work right, even though the other grep does work....
  #if [ "$(systemctl status $corsair_init | grep -c could\ not\ be\ found )" -eq 0 ]; then
  #  echo "Init system ready for initing";
  #fi


  # Prepare corsair user
  corsair_user="www-data"
  # and there's probably a better idea than assuming there's not already a user 33...
  if [ $(grep -c $corsair_user /etc/passwd) -eq 0 ]; then
    echo "Making $corsair_user user"
    $sudo_cmd groupadd -g 33 $corsair_user
    $sudo_cmd useradd \
      -g $corsair_user --no-user-group \
      --home-dir /var/www --no-create-home \
      --shell /usr/sbin/nologin \
      --system --uid 33 $corsair_user;
  fi



  # Prepare directories for corsair:
  if ! ls /etc/corsair >/dev/null 2>&1; then
    $sudo_cmd mkdir /etc/corsair
    $sudo_cmd chown -R root:$corsair_user /etc/corsair;
  fi
  if ! ls /etc/ssl/corsair >/dev/null 2>&1; then
    $sudo_cmd mkdir /etc/ssl/corsair
    $sudo_cmd chown -R $corsair_user:root /etc/ssl/corsair
    $sudo_cmd chmod 0770 /etc/ssl/corsair
  fi
  if ! ls /var/www >/dev/null 2>&1; then
    $sudo_cmd mkdir /var/www
    $sudo_cmd chown www-data:www-data /var/www
    $sudo_cmd chmod 555 /var/www
  fi

  # Make a dummy Corsairfile
  if ! ls /etc/corsair/Corsairfile >/dev/null 2>&1; then
    $sudo_cmd touch /etc/corsair/Corsairfile
    $sudo_cmd chown www-data:www-data /etc/corsair/Corsairfile
    $sudo_cmd chmod 444 /etc/corsair/Corsairfile
  fi

	# check installation
	$corsair_bin --version

  echo ""
	echo "Successfully installed"
  echo ""
  echo "Edit the Corsairfile at /etc/corsair/Corsairfile so that corsair can start"
  echo "See https://github.com/shipy4rd/corsair/docs/Corsairfile for more information"
  echo ""
  echo "To start corsair:"
  echo "sudo systemctl start corsair.service"
  echo ""
  echo "To enable automatic start on boot:"
  echo "sudo systemctl enable corsair.service"
  echo ""
  echo "A minimum ulimit of 8192 is suggested:"
  echo "sudo ulimit -n 8192"
  echo ""

	trap ERR
	return 0
}

install_corsair "$@"
