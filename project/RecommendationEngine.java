import java.io.*;
import java.util.*;
import java.util.stream.Collectors;

// This is the (complete) class that will generate the recommendations for a user
public class RecommendationEngine {
    private static final double R = 3.5; // Rating threshold for determining liked/disliked movies
    private static final int K = 10; // Minimum number of users that must like a movie for it to be considered
    private static final int N = 20; // Number of recommendations to return

    private final Map<Integer, Movie> movies; // Map to store movies by ID
    private final Map<Integer, User> users; // Map to store users by ID

    // constructs a recommendation engine from files
    public RecommendationEngine(String movieFile, String ratingFile) throws IOException {
        this.movies = new HashMap<>();
        this.users = new HashMap<>();
        readMovies(movieFile); // Read movie details from file
        readRatings(ratingFile); // Read ratings and user preferences from file
    }

    // Reads the Movie csv file of the MovieLens dataset
	// It populates the list of Movies
    private void readMovies(String filename) throws IOException {
        try (BufferedReader br = new BufferedReader(new FileReader(filename))) {
            br.readLine(); // Skip header
            String line;
            // Process each movie record in the file
            while ((line = br.readLine()) != null) {
                String[] parts = line.split(",(?=([^\"]*\"[^\"]*\")*[^\"]*$)", -1); // Split CSV with escaping quotes
                int id = Integer.parseInt(parts[0]); // Movie ID
                String title = parts[1].replace("\"", ""); // Movie title, remove quotes
                List<String> genres = Arrays.asList(parts[2].split("\\|")); // Movie genres split by '|'
                movies.put(id, new Movie(id, title, genres)); // Store movie object in the map
            }
        }
    }

    // Reads ratings data from a CSV file and updates user preferences
    private void readRatings(String filename) throws IOException {
        try (BufferedReader br = new BufferedReader(new FileReader(filename))) {
            br.readLine(); // Skip header
            String line;
            // Process each rating record in the file
            while ((line = br.readLine()) != null) {
                String[] parts = line.split(","); // Split CSV by commas
                int userId = Integer.parseInt(parts[0]); // User ID
                int movieId = Integer.parseInt(parts[1]); // Movie ID
                double rating = Double.parseDouble(parts[2]); // Rating value

                User user = users.computeIfAbsent(userId, User::new); // Create user if not already in map
                Movie movie = movies.get(movieId); // Get movie by ID
                
                if (rating >= R) { // If the rating is above the threshold, it's considered a liked movie
                    user.addLikedMovie(movieId);
                    movie.incrementLikedBy(); // Increment the count of users who liked the movie
                } else { // Otherwise, it's a disliked movie
                    user.addDislikedMovie(movieId);
                }
            }
        }
    }

    // Calculates the Jaccard similarity between two users based on their liked and disliked movies
    private double jaccardSimilarity(User u1, User u2) {
        // Get the intersection of liked movies for both users
        Set<Integer> commonLiked = new HashSet<>(u1.getLikedMovies());
        commonLiked.retainAll(u2.getLikedMovies());
        
        // Get the intersection of disliked movies for both users
        Set<Integer> commonDisliked = new HashSet<>(u1.getDislikedMovies());
        commonDisliked.retainAll(u2.getDislikedMovies());
        
        // Get the union of all rated movies for both users
        Set<Integer> allRated = new HashSet<>();
        allRated.addAll(u1.getLikedMovies());
        allRated.addAll(u1.getDislikedMovies());
        allRated.addAll(u2.getLikedMovies());
        allRated.addAll(u2.getDislikedMovies());

        if (allRated.isEmpty()) return 0.0; // If no movies have been rated, return similarity of 0.0
        // Return Jaccard similarity: ratio of common rated movies to total rated movies
        return (double)(commonLiked.size() + commonDisliked.size()) / allRated.size();
    }

    // Generates movie recommendations for a target user based on the Jaccard similarity of other users
    public List<Recommendation> generateRecommendations(int targetUserId) {
        User target = users.get(targetUserId); // Get the target user
        if (target == null) return Collections.emptyList(); // If the user doesn't exist, return an empty list

        return movies.values().parallelStream() // Process movies in parallel
            .filter(movie -> !target.hasRated(movie.getId())) // Exclude movies the target user has already rated
            .filter(movie -> movie.getLikedByCount() >= K) // Only consider movies liked by at least K users
            .map(movie -> {
                double score = 0.0;
                int likerCount = 0;
                
                // Calculate the recommendation score by comparing the target user to others who liked the movie
                for (User user : users.values()) {
                    if (user.getUserId() == targetUserId) continue; // Skip the target user
                    if (user.getLikedMovies().contains(movie.getId())) { // Check if the user liked the movie
                        score += jaccardSimilarity(target, user); // Add similarity to score
                        likerCount++;
                    }
                }
                
                // If there are users who liked the movie, create a recommendation object
                return likerCount > 0 
                    ? new Recommendation(target, movie, score / likerCount, likerCount) // Average similarity score
                    : null; // No recommendation if no users liked the movie
            })
            .filter(Objects::nonNull) // Remove null recommendations
            .sorted((a, b) -> Double.compare(b.getProbability(), a.getProbability())) // Sort by probability score in descending order
            .limit(N) // Limit to the top N recommendations
            .collect(Collectors.toList()); // Collect the results into a list
    }

    // Main method: entry point of the program, handles command line arguments
    public static void main(String[] args) {
        if (args.length != 3) { // Check if there are the correct number of arguments
            System.err.println("Usage: java RecommendationEngine <userID> <movies.csv> <ratings.csv>");
            System.exit(1); // Exit if arguments are incorrect
        }

        try {
            int targetUserId = Integer.parseInt(args[0]); // Get target user ID from the first argument
            long startTime = System.nanoTime(); // Record the start time

            RecommendationEngine engine = new RecommendationEngine(args[1], args[2]); // Create a RecommendationEngine object
            List<Recommendation> recommendations = engine.generateRecommendations(targetUserId); // Generate recommendations
            
            long endTime = System.nanoTime(); // Record the end time
            double executionTimeMs = (endTime - startTime) / 1_000_000.0; // Calculate execution time in milliseconds

            System.out.printf("Recommendations for user #  %d:%n", targetUserId); // Print recommendations header
            recommendations.forEach(rec -> System.out.printf(
                "%s at %.4f [%d]%n", // Print each movie's title, probability score, and number of users who liked it
                rec.getMovie().getTitle(),
                rec.getProbability(),
                rec.getNUsers()
            ));

            System.out.printf("%nExecution time: %.3fms%n", executionTimeMs); // Print execution time
        } catch (Exception e) {
            System.err.println("Error: " + e.getMessage()); // Print any errors
            e.printStackTrace(); // Print stack trace for debugging
        }
    }
}
