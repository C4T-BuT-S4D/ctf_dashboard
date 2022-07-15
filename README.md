[![Go Report Card](https://goreportcard.com/badge/github.com/C4T-BuT-S4D/ctf_dashboard)](https://goreportcard.com/report/github.com/C4T-BuT-S4D/ctf_dashboard)
[![release](https://github.com/C4T-BuT-S4D/ctf_dashboard/actions/workflows/release.yml/badge.svg)](https://github.com/C4T-BuT-S4D/ctf_dashboard/actions/workflows/release.yml)

# ctf_dashboard

Dashboard centralising all A&D CTF tools

To run, download the latest release and start the binary `dash`. Dashboard should be accessible at `http://0.0.0.0:8000`
.

Dashboard is designed to be integrated with [S4DFarm](https://github.com/C4T-BuT-S4D/S4DFarm),
[Neo](https://github.com/C4T-BuT-S4D/neo), [Goxy](https://github.com/pomo-mondreganto/goxy)
and MonGol (which is not open-source for now, but can be replaced by any external service using BasicAuth). The idea is
to set the same username/password on all services used by a CTF team and then access them from one centralized location
using links from the Dashboard
(links don't work in Safari though, and possibly other browsers, the project is tested in Chrome only). Dashboard also
allows uploading the players' ssh keys to vulnboxes and keeping track of the services for each vulnbox.

To use the ssh key upload feature, paste the private ssh key to the `resources/ssh_key` file and add the public part to
the authorized keys on the vulnbox.
