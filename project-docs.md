# London Underground Live Tracker

## Project Overview
Real-time London Underground visualization system showing live train positions across the network. Uses TfL API data with smart interpolation to create smooth, animated visualization of train movements.

## Why This Project?
- Demonstrates Go's concurrency strengths
- Real-time data processing
- Complex state management
- Interesting technical challenges
- Stands out from typical portfolio projects
- Actually useful application

## Technical Stack

### Backend (Vanilla Go)
- Native Go HTTP server
- WebSocket implementation
- Concurrent line tracking
- Rate limiting
- Error handling
- State management

### Frontend
- Next.js
- TypeScript
- WebSocket client
- D3.js for visualization
- Tailwind CSS

### External APIs
- TfL Unified API
  - 500 requests/min with key
  - 30-second update frequency
  - Arrival predictions
  - Service status

## Core Features

### MVP
1. Real-time train tracking
   - Position interpolation
   - Live arrival predictions
   - Smooth animations

2. Data Processing
   - Concurrent line polling
   - Efficient state management
   - Smart update batching

3. Visualization
   - Basic tube map
   - Train position updates
   - Service status
   - Simple animations

### Future Features
- Multiple line support
- Service disruptions
- Crowding data
- Journey planning
- Historical patterns

## Technical Challenges

### Data Processing
- Position interpolation from text descriptions
- Concurrent API polling
- Rate limit management
- State synchronization

### Real-time Updates
- WebSocket management
- Efficient client updates
- Connection handling
- Error recovery

### Visualization
- Smooth animation
- State management
- Performance optimization
- Responsive design

## MVP Development Plan

### Phase 1: Core Data Pipeline (3-4 days)
1. Basic Infrastructure
   - TfL API client
   - Rate limiting
   - Error handling
   - Basic logging

2. Line Management
   - Single line polling
   - Arrival predictions
   - Basic position estimation
   - State management

3. WebSocket Setup
   - Connection handling
   - Message protocol
   - State broadcasts

### Phase 2: Position Engine (2-3 days)
1. Train Position Estimation
   - Station coordinates
   - Basic interpolation
   - Speed calculations
   - Direction handling

2. State Management
   - Multiple trains
   - Update cycles
   - Data structures
   - Memory efficiency

### Phase 3: Basic Visualization (2-3 days)
1. Frontend Foundation
   - Next.js setup
   - WebSocket client
   - Basic tube map
   - Simple animations

2. Real-time Updates
   - Position updates
   - State sync
   - Error handling
   - Loading states

## TfL API Details

### Authentication
- Free registration required
- App ID and API Key needed
- Free tier available

### Rate Limits
- WITH key: 500 requests per min
- WITHOUT key: 50 requests per min
- No cost for standard API usage

### Data Format
```
Arrival Predictions:
- ISO timestamps with second precision
- Time to station in seconds
- Text-based current location
- Direction of travel
- Line ID and station name
```

### Update Frequency
- ~30 second refresh rate
- Real-time service status
- Text-based location updates

## Why Vanilla Go?
1. Better demonstration of Go knowledge
2. Cleaner WebSocket implementation
3. More control over concurrent design
4. Better for showing fundamentals
5. More impressive for portfolio
6. Core functionality doesn't need framework features

## Portfolio Impact
This project demonstrates:
1. Complex real-time system design
2. Concurrent data processing
3. Creative problem-solving
4. Working with API constraints
5. Full-stack development skills
6. Production-grade error handling

Stands out from typical portfolio projects by:
1. Solving a real problem
2. Using interesting algorithms
3. Handling complex state
4. Managing real-time data
5. Creating engaging visualization