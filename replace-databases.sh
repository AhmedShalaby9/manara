#!/bin/bash

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Replace Remote Databases with Local${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Database credentials
LOCAL_DB="manara"
LOCAL_USER="root"
LOCAL_PASS="2010961523"

STAGING_DB="manara_staging"
PRODUCTION_DB="manara_production"
SERVER_IP="75.119.144.109"

BACKUP_FILE="manara_backup_$(date +%Y%m%d_%H%M%S).sql"

# Step 1: Export local database
echo -e "${BLUE}📦 Exporting local database...${NC}"
mysqldump -u $LOCAL_USER -p$LOCAL_PASS $LOCAL_DB --single-transaction --set-gtid-purged=OFF > $BACKUP_FILE

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Local database exported: $BACKUP_FILE${NC}"
    FILE_SIZE=$(ls -lh $BACKUP_FILE | awk '{print $5}')
    echo -e "${GREEN}   File size: $FILE_SIZE${NC}"
else
    echo -e "${RED}❌ Failed to export local database${NC}"
    exit 1
fi
echo ""

# Step 2: Upload to server
echo -e "${BLUE}📤 Uploading to server...${NC}"
scp $BACKUP_FILE root@$SERVER_IP:/root/

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ File uploaded to server${NC}"
else
    echo -e "${RED}❌ Failed to upload file${NC}"
    exit 1
fi
echo ""

# Step 3: Replace staging database
echo -e "${YELLOW}⚠️  WARNING: This will DELETE all data in STAGING database!${NC}"
echo -e "${YELLOW}   Server: $SERVER_IP${NC}"
echo -e "${YELLOW}   Database: $STAGING_DB${NC}"
echo ""
read -p "Replace STAGING database? (yes/no): " CONFIRM_STAGING

if [ "$CONFIRM_STAGING" = "yes" ]; then
    echo -e "${BLUE}🔄 Replacing staging database...${NC}"
    ssh root@$SERVER_IP "mysql -u root $STAGING_DB < /root/$BACKUP_FILE"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Staging database replaced successfully${NC}"
    else
        echo -e "${RED}❌ Failed to replace staging database${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}⏭️  Skipped staging database${NC}"
fi
echo ""

# Step 4: Replace production database
echo -e "${RED}⚠️  DANGER: This will DELETE all data in PRODUCTION database!${NC}"
echo -e "${RED}   Server: $SERVER_IP${NC}"
echo -e "${RED}   Database: $PRODUCTION_DB${NC}"
echo -e "${RED}   This cannot be undone!${NC}"
echo ""
read -p "Replace PRODUCTION database? (TYPE 'yes' to confirm): " CONFIRM_PRODUCTION

if [ "$CONFIRM_PRODUCTION" = "yes" ]; then
    echo -e "${BLUE}🔄 Replacing production database...${NC}"
    ssh root@$SERVER_IP "mysql -u root $PRODUCTION_DB < /root/$BACKUP_FILE"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Production database replaced successfully${NC}"
    else
        echo -e "${RED}❌ Failed to replace production database${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}⏭️  Skipped production database${NC}"
fi
echo ""

# Step 5: Verify
echo -e "${BLUE}🔍 Verifying databases...${NC}"

echo -e "${BLUE}Staging tables:${NC}"
ssh root@$SERVER_IP "mysql -u root -e 'SHOW TABLES FROM $STAGING_DB;'"

echo ""
echo -e "${BLUE}Production tables:${NC}"
ssh root@$SERVER_IP "mysql -u root -e 'SHOW TABLES FROM $PRODUCTION_DB;'"

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}🎉 Database replacement completed!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "Summary:"
echo -e "  Backup file: $BACKUP_FILE"
echo -e "  Staging: $([ "$CONFIRM_STAGING" = "yes" ] && echo "✅ Replaced" || echo "⏭️  Skipped")"
echo -e "  Production: $([ "$CONFIRM_PRODUCTION" = "yes" ] && echo "✅ Replaced" || echo "⏭️  Skipped")"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo -e "  1. Restart staging: ssh root@$SERVER_IP 'systemctl restart manara-staging'"
echo -e "  2. Restart production: ssh root@$SERVER_IP 'systemctl restart manara-production'"
echo -e "  3. Test staging: curl https://staging.manaraco.net/api/roles"
echo -e "  4. Test production: curl https://api.manaraco.net/api/roles"
echo ""
echo -e "${GREEN}Backup file saved locally: $BACKUP_FILE${NC}"
echo ""
