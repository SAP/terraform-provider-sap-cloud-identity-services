### <u> OAuth2 Client Authentication </u>

You would require a valid **Client ID** and the corresponding **Client Secret** of a [System Administrator](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/add-administrators?version=Cloud#add-system-as-administrator) to get authenticated.

```NOTE: Refer to step 6 in the documentation linked above, section Secrets to fetch the required credentials. ```
 
There are multiple ways to configure your credentials.

1. You can configure them as part of the provider configuration as shown below:

    ```hcl
    provider "sci" {
        tenant_url = <your_tenant_url>
        client_id = <your_client_id>
        client_secret = <your_client_secret>
    }
    ```

2. You can export them as environment variables as shown below:

    #### Windows 
    
    If you use Windows CMD, do the export via the following commands:

    ```Shell
    set SCI_CLIENT_ID=<your_client_id>
    set SCI_CLIENT_SECRET=<your_client_secret>
    ```

    If you use Powershell, do the export via the following commands:

    ```Shell
    $Env:SCI_CLIENT_ID = '<your_client_id>'
    $Env:SCI_CLIENT_SECRET = '<your_client_secret>'
    ```

    #### Mac

    For Mac OS export the environment variables via:

    ```Shell
    export SCI_CLIENT_ID=<your_client_id>
    export SCI_CLIENT_SECRET=<your_client_secret>
    ```

    #### Linux

    For Linux export the environment variables via:

    ```Shell
    export SCI_CLIENT_ID=<your_client_id>
    export SCI_CLIENT_SECRET=<your_client_secret>
    ```