IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'e3-challenge')
BEGIN
    CREATE DATABASE [e3-challenge]
END
GO

IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'keycloak')
BEGIN
    CREATE DATABASE [keycloak]
END
GO


