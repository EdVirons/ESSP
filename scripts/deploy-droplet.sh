#!/bin/bash
set -e

# ESSP Deployment Script for Digital Ocean Droplet
# This script sets up Docker, pulls images, and deploys all services

# =============================================================================
# Configuration
# =============================================================================

DEPLOY_DIR="/opt/essp"
COMPOSE_FILE="docker-compose.prod.yml"
GITHUB_REGISTRY="ghcr.io"
IMAGE_TAG="${IMAGE_TAG:-latest}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# =============================================================================
# Install Docker if not present
# =============================================================================

install_docker() {
    if command -v docker &> /dev/null; then
        log_info "Docker is already installed: $(docker --version)"
        return 0
    fi

    log_info "Installing Docker..."

    # Update package index
    apt-get update

    # Install prerequisites
    apt-get install -y \
        ca-certificates \
        curl \
        gnupg \
        lsb-release

    # Add Docker's official GPG key
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    chmod a+r /etc/apt/keyrings/docker.gpg

    # Set up the Docker repository
    echo \
        "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
        $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

    # Install Docker Engine
    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    # Start and enable Docker
    systemctl start docker
    systemctl enable docker

    log_info "Docker installed successfully: $(docker --version)"
}

# =============================================================================
# Setup deployment directory
# =============================================================================

setup_deploy_dir() {
    log_info "Setting up deployment directory: ${DEPLOY_DIR}"

    mkdir -p "${DEPLOY_DIR}"
    mkdir -p "${DEPLOY_DIR}/certbot/conf"
    mkdir -p "${DEPLOY_DIR}/certbot/www"

    cd "${DEPLOY_DIR}"
}

# =============================================================================
# Login to GitHub Container Registry
# =============================================================================

docker_login() {
    if [ -z "${GHCR_TOKEN}" ]; then
        log_warn "GHCR_TOKEN not set, skipping registry login"
        return 0
    fi

    log_info "Logging into GitHub Container Registry..."
    echo "${GHCR_TOKEN}" | docker login "${GITHUB_REGISTRY}" -u "${GITHUB_USERNAME:-github}" --password-stdin
}

# =============================================================================
# Pull latest images
# =============================================================================

pull_images() {
    log_info "Pulling latest images with tag: ${IMAGE_TAG}"

    local repo="${GITHUB_REPOSITORY:-edvirons/essp}"

    docker pull "${GITHUB_REGISTRY}/${repo}/ims-api:${IMAGE_TAG}" || log_warn "Failed to pull ims-api"
    docker pull "${GITHUB_REGISTRY}/${repo}/ssot-school:${IMAGE_TAG}" || log_warn "Failed to pull ssot-school"
    docker pull "${GITHUB_REGISTRY}/${repo}/ssot-devices:${IMAGE_TAG}" || log_warn "Failed to pull ssot-devices"
    docker pull "${GITHUB_REGISTRY}/${repo}/ssot-parts:${IMAGE_TAG}" || log_warn "Failed to pull ssot-parts"
    docker pull "${GITHUB_REGISTRY}/${repo}/ssot-hr:${IMAGE_TAG}" || log_warn "Failed to pull ssot-hr"
    docker pull "${GITHUB_REGISTRY}/${repo}/sync-worker:${IMAGE_TAG}" || log_warn "Failed to pull sync-worker"

    log_info "Images pulled successfully"
}

# =============================================================================
# Deploy services
# =============================================================================

deploy_services() {
    log_info "Deploying services..."

    cd "${DEPLOY_DIR}"

    # Check if .env file exists
    if [ ! -f ".env" ]; then
        log_error ".env file not found in ${DEPLOY_DIR}"
        log_error "Please create .env file from .env.prod.example"
        exit 1
    fi

    # Stop existing services gracefully
    if docker compose -f "${COMPOSE_FILE}" ps -q 2>/dev/null | grep -q .; then
        log_info "Stopping existing services..."
        docker compose -f "${COMPOSE_FILE}" down --remove-orphans
    fi

    # Start services
    log_info "Starting services..."
    docker compose -f "${COMPOSE_FILE}" up -d

    log_info "Services deployed successfully"
}

# =============================================================================
# Run database migrations
# =============================================================================

run_migrations() {
    log_info "Waiting for database to be ready..."
    sleep 10

    log_info "Running database migrations..."

    # IMS API migrations
    docker compose -f "${COMPOSE_FILE}" exec -T ims-api /bin/sh -c "
        if [ -f /app/migrate ]; then
            /app/migrate -path /app/migrations -database \"\${PG_DSN}\" up
        fi
    " 2>/dev/null || log_warn "IMS API migrations skipped (no migrate binary or already up-to-date)"

    # SSOT migrations would run automatically on startup typically
    log_info "Migrations complete"
}

# =============================================================================
# Health check
# =============================================================================

health_check() {
    log_info "Running health checks..."

    local max_attempts=30
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        log_info "Health check attempt ${attempt}/${max_attempts}"

        # Check IMS API
        if curl -sf http://localhost/healthz > /dev/null 2>&1; then
            log_info "IMS API is healthy"

            # Show service status
            docker compose -f "${COMPOSE_FILE}" ps

            log_info "Deployment completed successfully!"
            return 0
        fi

        sleep 5
        attempt=$((attempt + 1))
    done

    log_error "Health check failed after ${max_attempts} attempts"
    docker compose -f "${COMPOSE_FILE}" logs --tail=50
    return 1
}

# =============================================================================
# Show logs
# =============================================================================

show_logs() {
    log_info "Showing recent logs..."
    docker compose -f "${COMPOSE_FILE}" logs --tail=100
}

# =============================================================================
# Cleanup old images
# =============================================================================

cleanup() {
    log_info "Cleaning up old Docker images..."
    docker image prune -af --filter "until=168h" || true
    docker system prune -f || true
}

# =============================================================================
# Main
# =============================================================================

main() {
    local command="${1:-deploy}"

    case "${command}" in
        install)
            install_docker
            ;;
        setup)
            setup_deploy_dir
            ;;
        pull)
            docker_login
            pull_images
            ;;
        deploy)
            install_docker
            setup_deploy_dir
            docker_login
            pull_images
            deploy_services
            run_migrations
            health_check
            cleanup
            ;;
        restart)
            cd "${DEPLOY_DIR}"
            docker compose -f "${COMPOSE_FILE}" restart
            health_check
            ;;
        stop)
            cd "${DEPLOY_DIR}"
            docker compose -f "${COMPOSE_FILE}" down
            ;;
        logs)
            cd "${DEPLOY_DIR}"
            show_logs
            ;;
        status)
            cd "${DEPLOY_DIR}"
            docker compose -f "${COMPOSE_FILE}" ps
            ;;
        health)
            health_check
            ;;
        cleanup)
            cleanup
            ;;
        *)
            echo "Usage: $0 {install|setup|pull|deploy|restart|stop|logs|status|health|cleanup}"
            exit 1
            ;;
    esac
}

main "$@"
