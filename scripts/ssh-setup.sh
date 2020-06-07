#!/usr/bin/env bash
openssl aes-256-cbc -K $encrypted_da7eec2e51b3_key -iv $encrypted_da7eec2e51b3_iv -in travis_rsa.enc -out /tmp/travis_rsa -d
chmod 400 /tmp/travis_rsa
eval "$(ssh-agent -s)"
ssh-add /tmp/travis_rsa
