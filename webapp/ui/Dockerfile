# Use the official Node.js base image
FROM node:16-alpine

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json files into the working directory
COPY . .

# Install dependencies
RUN npm i

# Expose the port on which the app will run
EXPOSE 3000

# Start the application
CMD ["npm", "start"]