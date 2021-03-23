FROM node:latest

RUN npm install -g truffle --loglevel verbose
RUN mkdir -p openzeppelin
WORKDIR /openzeppelin
COPY ./openzeppelin-contracts/package.json /openzeppelin
COPY ./openzeppelin-contracts/hardhat /openzeppelin/hardhat
COPY ./openzeppelin-contracts/hardhat.config.js /openzeppelin
RUN yarn install
COPY ./openzeppelin-contracts/contracts /openzeppelin/contracts
COPY ./truffle-config.js /openzeppelin
RUN truffle compile
COPY ./openzeppelin-contracts/scripts /openzeppelin/scripts
COPY ./openzeppelin-contracts/ /openzeppelin
COPY ./truffle-config.js /openzeppelin

CMD [ "truffle", "test" ]
