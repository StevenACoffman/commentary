# in order to download release artifacts from github, you have to first retrieve the
# list of asset URLs using the github repo REST API. Use the asset URL to download 
# the artifact as a octet-stream data stream. You will need to get an access token 
# from "settings -> developer settings -> personal access tokens" on the github UI
#!/bin/bash -e

tag="$(git tag -l --sort=-version:refname v* | head -1)"
tag_without_v="$(echo -n "${tag}"| cut -c2-)"
# list_asset_url="https://api.github.com/repos/${GITHUB_REPOSITORY}/releases/tags/${tag}"
# RUNNER_OS is Linux, Windows, or macOS
# RUNNER_ARCH	is X86, X64, ARM, or ARM64.
# feels like too much work to construct $artifact to be commentary_0.6.0_linux_amd64.tar.gz
# and then get url for artifact with name==$artifact
# asset_url=$(curl "${list_asset_url}" -H "Authorization: bearer ${GITHUB_TOKEN}" | jq ".assets[] | select(.name==\"${artifact}\") | .url" | sed 's/\"//g')

asset_url="https://github.com/${GITHUB_REPOSITORY}/releases/download/${tag}/commentary_${tag_without_v}_darwin_amd64.tar.gz"

# download the artifact
#curl -vLJO -H "Authorization: bearer ${GITHUB_TOKEN}" -H 'Accept: application/octet-stream' \
#     "${asset_url}" | tar xvzf - commentary
# tar xvzf "commentary_${tag_without_v}_darwin_amd64.tar.gz" commentary
curl -LJ -H "Authorization: bearer ${GITHUB_TOKEN}" -H 'Accept: application/octet-stream' \
     "${asset_url}" | tar xvzf - commentary
./commentary