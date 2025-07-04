### <u> Basic Authentication </u>

You would require a valid **username** and **password** of a [User Administrator](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/add-administrators?version=Cloud#add-user-as-administrator) to get authenticated.
 
There's multiple ways to configure your credentials.

1. You can configure them as part of the provider configuration as shown below:

    ```hcl
    provider "sci" {
        tenant_url = <your_tenant_url>
        username = <your_username>
        password = <your_password>
    }
    ```

2. You can export them as environment variables as shown below:

    #### Windows 
    
    If you use Windows CMD, do the export via the following commands:

    ```Shell
    set SCI_USERNAME=<your_username>
    set SCI_PASSWORD=<your_password>
    ```

    If you use Powershell, do the export via the following commands:

    ```Shell
    $Env:SCI_USERNAME = '<your_username>'
    $Env:SCI_PASSWORD = '<your_password>'
    ```

    #### Mac

    For Mac OS export the environment variables via:

    ```Shell
    export SCI_USERNAME=<your_username>
    export SCI_PASSWORD=<your_password>
    ```

    #### Linux

    For Linux export the environment variables via:

    ```Shell
    export SCI_USERNAME=<your_username>
    export SCI_PASSWORD=<your_password>
    ```