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

1. **Clone the repository:**

   ```bash
   git clone https://github.com/neilfarmer/go-lazy-openstack.git
   cd go-lazy-openstack
   ```

2. **Build:**

   ```bash
   go build -o go-lazy-openstack
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
   ./go-lazy-openstack
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

## Acknowledgments

- [tview](https://github.com/rivo/tview)
- [tcell](https://github.com/gdamore/tcell)

---

*Happy OpenStacking!*
