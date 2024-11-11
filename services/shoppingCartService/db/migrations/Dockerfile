FROM gomicro/goose

# Add migration files
ADD /*.sql /migrations/

# Set entrypoint
ENTRYPOINT ["sh", "-c", "goose -dir=/migrations postgres \"host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=$DB_SSLMODE\" up"]
