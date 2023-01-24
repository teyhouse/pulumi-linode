package main

import (
	"linode/ansb"
	"strconv"

	"github.com/pulumi/pulumi-linode/sdk/v3/go/linode"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const instancecount = 2

func main() {
	ansb.CreateAnsibleHostTemplate()

	pulumi.Run(func(ctx *pulumi.Context) error {

		//Create Firewall
		myFirewall, err := linode.NewFirewall(ctx, "myFirewall", &linode.FirewallArgs{
			Label: pulumi.String("my_firewall"),
			Inbounds: linode.FirewallInboundArray{
				&linode.FirewallInboundArgs{
					Label:    pulumi.String("http"),
					Action:   pulumi.String("ACCEPT"),
					Protocol: pulumi.String("TCP"),
					Ports:    pulumi.String("80"),
					Ipv4s: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
					},
					Ipv6s: pulumi.StringArray{
						pulumi.String("::/0"),
					},
				},
				&linode.FirewallInboundArgs{
					Label:    pulumi.String("ALLOW-SSH-TT"),
					Action:   pulumi.String("ACCEPT"),
					Protocol: pulumi.String("TCP"),
					Ports:    pulumi.String("22"),
					Ipv4s: pulumi.StringArray{
						pulumi.String("77.180.139.68/32"),
					},
				},
			},
			InboundPolicy:  pulumi.String("DROP"),
			OutboundPolicy: pulumi.String("ACCEPT"),
		})
		if err != nil {
			return err
		}

		//Collector for all instances
		instancelist := []linode.Instance{}

		//Create amount of instances defined in instancecount
		for i := 0; i < instancecount; i++ {
			myInstance, err := linode.NewInstance(ctx, "myInstance"+strconv.Itoa(i), &linode.InstanceArgs{
				Label:    pulumi.String("my_instance-" + strconv.Itoa(i)),
				Image:    pulumi.String("linode/ubuntu22.10"),
				Region:   pulumi.String("eu-central"),
				Type:     pulumi.String("g6-nanode-1"),
				RootPass: pulumi.String("S3Crb0GsFamOri!"),
				AuthorizedKeys: pulumi.StringArray{
					pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDYu3/TTJs3iDcQbkstLkE8oFeIzIm9LNTz2n6hkziJgWSf8YpysTCMWSm9yJCvswq1oiCQiKZGUMycaJ6qK4Dlls79N0yhmQwav8MwhQiCmvSPJGr4+bakUJvLw9l4XXz+/oJ6QGBAi0bnybBxarmkm1iMj9DqM8Q58fwSWH9U6/PVtVpWDEksNORcMIHecBIsEwnOA+/BV5ZiwUmg5W0qHUUowS9LjWFPFlejCJbcHsSAuc/CC+OLw3rwhBq5ti9cdXToGTtwt7aMlySr7C/jEopDtpaew8sg7DR7zBaDv2ZV0PnD+SfKTQuz+rBbjDlA7I/gHb8nKRdasoW3xx7Z9KpWTbFbQz3S4PuPPgMYMtvLHGZvO9fCrOu8dttkM4nfTLvveHn23YDKpxn4icK1SpSg7O9d3b8QtLtVkocRljEOs4LFerJuC8wAMm06xi8jArXtbqVKhPfy4lCA1uw50PCQChDGsalpM+ce7WSMP4Go4SrSDL36VnKs734GAqk= administrator@ubu"),
					pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCww0PWD4UhNK7SzvhCOouHJlwJ5nSL/br1GEuFcMlsgiYw+CDzXZN3Q3Hqwq7iiQU/fxX6Dcjm0SOh1QfuCzsJsqX23KHW+USjlSW3kM2W4I4cQadkw5te3d4qa3ydjDJ0EMg4nZk6skmhVEBpgloup58W9yr5sBkLBLMdWllW/WTieJ0cM5ipym4RFaomeZ1CS8UZNThpYlnuQzCKhx6kxFxzDssbYGBaZCNfrl2hK9r8dW7SOPkH7C7UnpbZn0xpfexuAJpPJq+rbjBA5i7RjNTVxFBL7tM0dQKrNQSt0uHvkdfooP3IYq72ILPCk6vOASKkER02BkJxGZwSA/vMUURb5S2tipH6VEoUVRGeX+7JYC9opdWYjgYvlLdo1jZfK/h+Am15TbK/bcgsqZEXBOfGjsIoSES5nzesfgGCvRZt6UePTORHBptwE41DhKNPrG1BuiDVJ4zbUo8N2b1r0i09zIpbk+WtkQGjGSMzF10YrEVfKxP6j7TNP8Qf5H8= teyhouse@TEYHOUSE-PC"),
				},
				SwapSize: pulumi.Int(256),
				Interfaces: linode.InstanceInterfaceArray{
					//eth0
					linode.InstanceInterfaceArgs{
						Purpose: pulumi.String("public"),
					},
					//eth1
					linode.InstanceInterfaceArgs{
						IpamAddress: pulumi.String("192.168.1." + strconv.Itoa(i+1) + "/24"), //Silly fix so we don't use zero (which is the network address) - this must be improved
						Label:       pulumi.String("1"),
						Purpose:     pulumi.String("vlan"),
					},
				},
			})
			if err != nil {
				return err
			}
			instancelist = append(instancelist, *myInstance)
		}

		//Attach Firewall to VM
		for i, instance := range instancelist {
			var conversionCallback = func(val string) (int, error) {
				return strconv.Atoi(val)
			}

			_, err = linode.NewFirewallDevice(ctx, "myDevice-"+strconv.Itoa(i), &linode.FirewallDeviceArgs{
				FirewallId: myFirewall.ID().ToStringOutput().ApplyT(conversionCallback).(pulumi.IntOutput),
				EntityId:   instance.ID().ToStringOutput().ApplyT(conversionCallback).(pulumi.IntOutput),
			}, pulumi.Parent(myFirewall), pulumi.ReplaceOnChanges([]string{"*"}))
			if err != nil {
				return err
			}
		}

		//Append to Ansible Hosts-List
		for _, instance := range instancelist {
			_ = instance.IpAddress.ApplyT(func(IpAddress string) string {
				ansb.AppendAnsible(IpAddress)
				return IpAddress
			}).(pulumi.StringOutput)
		}

		return nil
	})

	//Append Ansible Group-Identifier & Run Ansible-Playbook
	ansb.AppendAnsible("[instances]")
	ansb.RunPlayBook()
}
