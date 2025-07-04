### <u> Certificate Based Authentication </u>

You would require a valid **p12 certificate** and the corresponding **password** of a [System Administrator](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/add-administrators?version=Cloud#add-system-as-administrator) to get authenticated.
 
```NOTE: Refer to step 6 in the documentation linked above, section Certificates to fetch the required credentials. ```

You can only configure the credentials as part of the provider configuration as shown below:

 ```hcl
provider "sci" {
    tenant_url = <your_tenant_url>
    p12_certificate_content = <your_p12_certificate>
    p12_certificate_password = <your_p12_certificate_password>
}
```

Ensure to paste the ***content*** of your p12 certificate rather than the the ***file path***.
You can even use the function `filebase64("path_to_certificate.p12")` to load the file content. 

