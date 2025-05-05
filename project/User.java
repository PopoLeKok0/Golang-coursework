import java.util.HashSet;
import java.util.Set;

public class User {
    private final int userId;
    private final Set<Integer> likedMovies;
    private final Set<Integer> dislikedMovies;

    // Constructor
    public User(int userId) {
        this.userId = userId;
        this.likedMovies = new HashSet<>();
        this.dislikedMovies = new HashSet<>();
    }

    // gets the user id
    public int getUserId() { 

        return userId; 
    }

    // gets the liked movies
    public Set<Integer> getLikedMovies() { 

        return likedMovies; 
    }

    // gets the disliked movies
    public Set<Integer> getDislikedMovies() { 
        return dislikedMovies; 
    }

    // Add a movie to the liked list
    public void addLikedMovie(int movieId) { 

        likedMovies.add(movieId); 
    }

    // Add a movie to the disliked list
    public void addDislikedMovie(int movieId) { 
        
        dislikedMovies.add(movieId);
    }

    // Check if the user has rated a movie
    public boolean hasRated(int movieId) { 
        return likedMovies.contains(movieId) || dislikedMovies.contains(movieId); 
    }
}