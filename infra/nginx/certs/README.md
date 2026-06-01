# Self-signed dev cert (replace in production with Let's Encrypt)
# Generate:
#   openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
#     -keyout privkey.pem -out fullchain.pem -subj "/CN=localhost"
