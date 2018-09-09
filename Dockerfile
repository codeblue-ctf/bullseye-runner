FROM node:10

WORKDIR /usr/local/bulls-eye-runner
ADD ./package.json .
RUN npm i

ADD . /usr/local/bulls-eye-runner
CMD ["npm", "run"]