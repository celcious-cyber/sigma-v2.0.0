$loginBody = @{
    nik = "1234567890123456" # Default admin NIK or we can try default credential
    password = "password" # Or we can just bypass
} | ConvertTo-Json

# Actually let's just make a Go script that hits the database and does json serialization, wait we already did that!
# Let's try to query the running API using a scratch Go script instead of PowerShell, it's easier to handle HTTP requests.
