# Simple frontend image - just install dependencies
FROM node:20-alpine

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Expose port
EXPOSE 3000


RUN  npm build

# Start development server
CMD ["npm", "run", "start"]