#! /usr/bin/bash
#Must be run as sudo

#Copy the scripts files to the postgres directory
cp -r ../database/ ~postgres/
cd ~postgres/

#Run scripts with psql: db_superbchat.sql  account.sql  superchat.sql
echo "Creating 'superbchat' database"
sudo -u postgres -H -- psql -U postgres -d postgres -f database/db_superbchat.sql
echo "Creating 'account' table"
sudo -u postgres -H -- psql -U postgres -d superbchat -f database/account.sql
echo "Creating 'superchat' table"
sudo -u postgres -H -- psql -U postgres -d superbchat -f database/superchat.sql
#Create 'web' user
echo "Creating 'web' user"
sudo -u postgres -H -- psql -U postgres -d superbchat -c "CREATE USER web WITH PASSWORD 'pass';"
sudo -u postgres -H -- psql -U postgres -d superbchat -c "GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO web;"
sudo -u postgres -H -- psql -U postgres -d superbchat -c "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO web;"

#Remove scripts from postgres directory
rm -r -f database/

