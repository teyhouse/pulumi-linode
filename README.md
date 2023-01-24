# 💡  Pulumi Linode Go Example (with Ansible)
This is meant as an example on how to create Linode-Instances (VMs) and configure them using Ansible.
The Ansible-Hosts is auto-generated on every run. This project was created based on the official Pulumi-Go Template.  
  
# 📃 Requirements
- .Go >= 1.18  
https://go.dev/doc/install 
- Pulumi:  
https://www.pulumi.com/docs/get-started/install/
- Pulumi Linode Provider:  
```go get github.com/pulumi/pulumi-linode/sdk/v3```  
(check for newer version on: github.com/pulumi/pulumi-linode/sdk/v3/go/linode)  
- Linode API Token - create on this page:  
https://cloud.linode.com/profile/tokens
  
- Make sure to set your Linode-Token as Pulumi-Secret:
``pulumi config set linode:token XXXXXXXXXXXXXX --secret``
  
# 🚫 Limitations
This example is meant as a showcase, so certain aspects have been purposely simplified.  
In real environments, you should create more abstraction and therefore usability for your resources (go-modules for every resource type, for example).  
Certain configuration elements (resource amounts, ssh keys, etc.) should be moved to configuration files (yaml, json, or whatever you prefer).  
If you are planning on using VLANs, please implement a real IP generation algorithm; the current implementation is just thin duct tape.
  
# 📖 API
https://www.pulumi.com/registry/packages/linode/api-docs/  
  
# 🛠 How to run 
Preview changes:  
```pulumi preview```
  
Run:  
```pulumi up```
  
Revert changes:  
```pulumi destroy```

![screenshot](pulumi.png?raw=true)
