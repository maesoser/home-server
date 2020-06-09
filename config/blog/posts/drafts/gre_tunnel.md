# Linux GRE Tunnel

So the idea is to try GRE tunnels:



```
 
```



```
auto gre1
iface gre1 inet static
address 10.10.10.10  # inside side a address
netmask 255.255.255.0
pre-up ip tunnel add gre1 mode gre remote 173.230.145.76 local 173.230.147.224
post-down ip tunnel del gre1

pre-up ip route add 173.230.145.76/32 via (ip -4 route list default | awk '{print $3}' )

to avoid tunnel Loops. Tunnel Loop can appear if, in your case 173.230.145.0/24 would be learnd via gre1
alex-eri commented on 9 Sep 2019
ip route get ${REMOTE} | awk '{print $3}'

```





https://developers.redhat.com/blog/2019/05/17/an-introduction-to-linux-virtual-interfaces-tunnels/#sit



## References

https://www.tldp.org/HOWTO/Adv-Routing-HOWTO/lartc.tunnel.gre.html 

https://www.krijnders.net/index.php?page=personal/interfaces



SSH tunnel:

https://www.tummy.com/blogs/2013/03/04/ssh-network-tunneling/



# IPFS cluster

https://gist.github.com/StephanieSunshine/92a3af3fc8577103906ff8142a4349f5

