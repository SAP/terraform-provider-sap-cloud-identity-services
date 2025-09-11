### <u> X.509 Certificate Authentication </u>

You would require a valid **p12 certificate** and the corresponding **password** of a [System Administrator](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/add-administrators?version=Cloud#add-system-as-administrator) to get authenticated.
 
```NOTE: Refer to step 6 in the documentation linked above, section Certificates to fetch the required credentials. ```

You can configure the credentials as part of the provider configuration as shown below:

 ```hcl
provider "sci" {
    tenant_url = <your_tenant_url>
    p12_certificate_content = <your_p12_certificate>
    p12_certificate_password = <your_p12_certificate_password>
}
```

Ensure to paste the ***content*** of your p12 certificate rather than the ***file path***.
You can even use the function `filebase64("path_to_certificate.p12")` to load the file content. 

2. You can export the Certificate Password as an environment variable as shown below:

    #### Windows 

    If you use Windows CMD, do the export via the following commands:

    ```Shell
    set SCI_P12_CERTIFICATE_PASSWORD=<your_password>
    ```

    If you use Powershell, do the export via the following commands:

    ```Shell
    $Env:SCI_P12_CERTIFICATE_PASSWORD = '<your_password>'
    ```

    #### Mac

    For Mac OS export the environment variable via:

    ```Shell
    export SCI_P12_CERTIFICATE_PASSWORD=<your_password>
    ```

    #### Linux

    For Linux export the environment variable via:

    ```Shell
    export SCI_P12_CERTIFICATE_PASSWORD=<your_password>
    ```

    **The P12 Certificate itself would still have to be configured as a schema parameter even if the password is exported as an environment variable**.