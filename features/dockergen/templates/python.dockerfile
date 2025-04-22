# Use a slim or Alpine base image for smaller size
FROM python:3.12-alpine

# Set working directory
WORKDIR /app

# Copy only requirements first for better caching
COPY requirements.txt .

# Install dependencies
RUN pip install --upgrade pip && pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Document the port (optional, for web apps)
EXPOSE 8000

# Start the application (adjust as needed)
CMD ["python", "app.py"]
