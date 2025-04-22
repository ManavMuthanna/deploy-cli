# Use a lightweight, explicit base image
FROM node:22-alpine

# Set working directory
WORKDIR /usr/src/app

# Copy dependency descriptors first for better caching
COPY package*.json ./

# Install only production dependencies deterministically
RUN npm ci --only=production

# Copy application source code
COPY . .

# Set environment variable for production optimizations
ENV NODE_ENV=production

# Document the port your app runs on
EXPOSE 3000

# Start the application
CMD ["node", "app.js"]
