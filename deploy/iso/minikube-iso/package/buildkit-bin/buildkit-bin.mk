################################################################################
#
# buildkit-bin
#
################################################################################

BUILDKIT_BIN_VERSION = v0.9.0
BUILDKIT_BIN_COMMIT = c8bb937807d405d92be91f06ce2629e6202ac7a9
BUILDKIT_BIN_SITE = https://github.com/moby/buildkit/releases/download/$(BUILDKIT_BIN_VERSION)
BUILDKIT_BIN_SOURCE = buildkit-$(BUILDKIT_BIN_VERSION).linux-amd64.tar.gz

# https://github.com/opencontainers/runc.git
BUILDKIT_RUNC_VERSION = 12644e614e25b05da6fd08a38ffa0cfe1903fdec

define BUILDKIT_BIN_INSTALL_TARGET_CMDS
	$(INSTALL) -D -m 0755 \
		$(@D)/buildctl \
		$(TARGET_DIR)/usr/bin
	$(INSTALL) -D -m 0755 \
		$(@D)/buildkit-runc \
		$(TARGET_DIR)/usr/sbin
	$(INSTALL) -D -m 0755 \
		$(@D)/buildkit-qemu-* \
		$(TARGET_DIR)/usr/sbin
	$(INSTALL) -D -m 0755 \
		$(@D)/buildkitd \
		$(TARGET_DIR)/usr/sbin
endef

$(eval $(generic-package))
