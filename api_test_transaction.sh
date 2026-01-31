#!/bin/bash

# Base URL
BASE_URL="http://localhost:8080"

# 1. Login as Budi to get Token
echo "1. Logging in as Budi (User ID 2)..."
# Password from seed.go: pAsswj@1873
LOGIN_RESP=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "budi@mail.com",
    "password": "pAsswj@1873"
  }')

# Extract Token (Simple parsing, requires jq installed, or manual copy)
TOKEN=$(echo $LOGIN_RESP | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "Failed to login. Response: $LOGIN_RESP"
    exit 1
fi
echo "✅ Token collected."

# 2. Create Transaction (Success - Within Limit)
# Budi has 100,000 limit for Tenor 1
echo -e "\n2. Creating Transaction (Valid - OTR 50,000 < 100,000 Limit)..."
curl -i -X POST "$BASE_URL/api/transaction/" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "contract_number": "CTR-BUDI-001",
    "otr": 50000,
    "admin_fee": 1000,
    "installment_amount": 51000,
    "interest_amount": 1000,
    "asset_name": "Rice Cooker",
    "tenor": 1
  }'

# 3. Create Transaction (Failure - Exceeds Limit)
echo -e "\n\n3. Creating Transaction (Invalid - OTR 150,000 > 100,000 Limit)..."
curl -i -X POST "$BASE_URL/api/transaction/" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "contract_number": "CTR-BUDI-002",
    "otr": 150000,
    "admin_fee": 2000,
    "installment_amount": 152000,
    "interest_amount": 2000,
    "asset_name": "Smartphone",
    "tenor": 1
  }'
