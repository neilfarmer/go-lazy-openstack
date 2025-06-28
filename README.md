# go-lazy-openstack

A lightweight, interactive terminal UI for quickly navigating and managing OpenStack resources.  
Built with [tview](https://github.com/rivo/tview) and [tcell](https://github.com/gdamore/tcell), this tool lets you browse servers, images, flavors, projects, volumes, and networks across OpenStack projects from the comfort of your terminal.

---

## Features

- **Project Quick Switch:** Instantly change the active OpenStack project.
- **Resource Browsing:** View servers, images, flavors, volumes, and networks in an interactive list.
- **Details Pane:** See detailed information for each selected resource.
- **Keyboard-Driven:** Fast navigation using intuitive key bindings.
- **Command Prompt:** Type commands or use shortcuts for navigation.

---

## Installation

1. Run install commands

```
curl -LO https://github.com/neilfarmer/go-lazy-openstack/releases/download/v0.1.0/go-lazy-openstack-darwin-arm64
chmod +x go-lazy-openstack-darwin-arm64
sudo mv go-lazy-openstack-darwin-arm64 /usr/local/bin/go-lazy-openstack
```

---

## Usage

1. **Set your OpenStack environment variables.**  
   This app expects the usual OpenStack CLI variables, such as:

   ```bash
   export OS_AUTH_URL=https://openstack.example.com:5000/v3
   export OS_PROJECT_NAME=myproject
   export OS_USERNAME=myuser
   export OS_PASSWORD=mypassword
   export OS_USER_DOMAIN_NAME=Default
   export OS_PROJECT_DOMAIN_NAME=Default
   ```

2. **Run the application:**

   ```bash
   go-lazy-openstack
   ```

3. **Navigation:**

   - Use the following single-key shortcuts to access resources:
     - `p` — Projects
     - `i` — Images
     - `f` — Flavors
     - `s` — Servers
     - `n` — Networks
     - `v` — Volumes
     - `q` — Quit
   - Use arrow keys and Enter to select items.
   - Press `:` to open the command prompt for typing resource names (e.g., `servers`, `images`).

4. **Switch Projects:**
   - Go to Projects (`p`), select a project, and your selection is set as the active OpenStack project for the session.

---

## Requirements

- Go 1.18+
- OpenStack API access and environment variables set

---

## Contributing

Pull requests and issues are welcome! Please open an issue to discuss major changes.

---

## Releases

1. Create a git tag with a semantic version

```
git tag v0.1.0
git push origin v0.1.0
```

2. This will trigger the Github Action workflow to:

- Build the binary for:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64)

---

## Acknowledgments

- [tview](https://github.com/rivo/tview)
- [tcell](https://github.com/gdamore/tcell)

---

*Happy OpenStacking!*
