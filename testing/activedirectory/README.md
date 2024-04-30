## Test ds-load ldap plugin with active directory

1. Install active directory in windows server 2022 https://computingforgeeks.com/install-active-directory-domain-services-in-windows-server/

2. Configure SSH to window server https://computingforgeeks.com/configure-openssh-server-on-windows-server/

3. Forward the ldap port to your local machine
e.g.
```
ssh  -L 1389:127.0.0.1:389 asert@123.123.123.123
```

4. Use the sample config and add the necessary credentials to connect to ldap.