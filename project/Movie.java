import java.util.List;

public class Movie {
    private final int movieID;
    private final String title;
    private final List<String> genres;
    private int likedByCount;

    // constructs a movie
    public Movie(int movieID, String title, List<String> genres) {
        this.movieID = movieID;
        this.title = title;
        this.genres = genres;
        this.likedByCount = 0;
    }
    // gets the ID
    public int getId() { 

        return movieID; 
    }
    // get the movie title
    public String getTitle() { 

        return title; 
    }
    
    // get the movie genres
    public List<String> getGenres() { 

        return genres; 
    }

    // get the number of likes
    public int getLikedByCount() { 
        
        return likedByCount; 
    }

    // increment the number of likes
    public void incrementLikedBy() { 

        likedByCount++; 
    }
}