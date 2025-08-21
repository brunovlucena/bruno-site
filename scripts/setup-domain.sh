#!/bin/bash

# Bruno Site Domain Setup Script
# This script helps set up the domain and verify nginx-ingress configuration

set -e

echo "🌐 Bruno Site Domain Setup"
echo "=========================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if nginx-ingress is running
echo -e "${YELLOW}📋 Checking nginx-ingress status...${NC}"
if kubectl get pods -n nginx-ingress | grep -q "Running"; then
    echo -e "${GREEN}✅ nginx-ingress is running${NC}"
else
    echo -e "${RED}❌ nginx-ingress is not running${NC}"
    exit 1
fi

# Get the external IP of nginx-ingress
echo -e "${YELLOW}🌍 Getting nginx-ingress external IP...${NC}"
EXTERNAL_IP=$(kubectl get svc -n nginx-ingress nginx-ingress-ingress-nginx-controller -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "Not available")

if [ "$EXTERNAL_IP" = "Not available" ]; then
    echo -e "${YELLOW}⚠️  External IP not available. This might be a local cluster.${NC}"
    echo -e "${YELLOW}💡 For local development, you can use port-forwarding:${NC}"
    echo -e "${YELLOW}   kubectl port-forward -n nginx-ingress svc/nginx-ingress-ingress-nginx-controller 80:80 443:443${NC}"
else
    echo -e "${GREEN}✅ External IP: $EXTERNAL_IP${NC}"
    echo -e "${YELLOW}💡 Configure your DNS to point lucena.cloud to: $EXTERNAL_IP${NC}"
fi

# Check if ClusterIssuer exists
echo -e "${YELLOW}🔐 Checking Let's Encrypt ClusterIssuer...${NC}"
if kubectl get clusterissuer letsencrypt-prod >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Let's Encrypt ClusterIssuer exists${NC}"
else
    echo -e "${YELLOW}📋 Creating Let's Encrypt ClusterIssuer...${NC}"
    kubectl apply -f chart/templates/cluster-issuer.yaml
    echo -e "${GREEN}✅ ClusterIssuer created${NC}"
fi

# Check ingress status
echo -e "${YELLOW}🚪 Checking ingress status...${NC}"
if kubectl get ingress -n bruno >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Ingress exists in bruno namespace${NC}"
    kubectl get ingress -n bruno
else
    echo -e "${RED}❌ No ingress found in bruno namespace${NC}"
fi

# Check certificate status
echo -e "${YELLOW}🔒 Checking certificate status...${NC}"
if kubectl get certificate -n bruno >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Certificates found:${NC}"
    kubectl get certificate -n bruno
else
    echo -e "${YELLOW}📋 No certificates found yet. They will be created automatically when the ingress is accessed.${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Setup complete!${NC}"
echo ""
echo -e "${YELLOW}📋 Next steps:${NC}"
echo "1. Configure your DNS to point lucena.cloud to your cluster's external IP"
echo "2. Wait for the certificate to be issued (can take a few minutes)"
echo "3. Access your site at https://lucena.cloud"
echo ""
echo -e "${YELLOW}🔍 To monitor certificate status:${NC}"
echo "   kubectl get certificate -n bruno -w"
echo ""
echo -e "${YELLOW}🔍 To check ingress status:${NC}"
echo "   kubectl describe ingress -n bruno"
