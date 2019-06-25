package me.photolab;

import java.io.File;
import java.util.ArrayList;
import me.photolab.ClientPhotolab;

public class App
{
    public static void main( String[] args )
    {
        ClientPhotolab cl = new ClientPhotolab();
        try {
            String resoursesFilename = "resources.zip";
            File f = new File(resoursesFilename);
            if (!f.exists()) {
                cl.DownloadFile("http://soft.photolab.me/samples/resources.zip", resoursesFilename);
            }

            String contentFilename = "girl.jpg";
            File fIm = new File(contentFilename);
            if (!fIm.exists()) {
                cl.DownloadFile("http://soft.photolab.me/samples/girl.jpg", contentFilename);
            }

            String contentUrl = cl.ImageUpload(contentFilename);
            System.out.println(String.format("content url: %s", contentUrl));

            String templateName = cl.TemplateUpload(resoursesFilename);
            System.out.println(String.format("template name: %s", templateName));

            ArrayList<ImageRequest> testArr = new ArrayList<ImageRequest>();
            ImageRequest testIm = new ImageRequest();
            testIm.url = contentUrl;
            testIm.crop = "0,0,1,1";
            testIm.flip = 0;
            testIm.rotate = 0;
            testArr.add(testIm);
            ImageRequest[] testArrRequests = new ImageRequest[testArr.size()];
            testArrRequests = testArr.toArray(testArrRequests);

            String firstResultUrl = cl.TemplateProcess(templateName, testArrRequests);
            System.out.println(String.format("--- for template_name: %s, result_url: %s", templateName, firstResultUrl));

            long[] comboIds = new long[]{5635874, 3124589};
            for (int i = 0; i < comboIds.length; i++) {
                long comboId = comboIds[i];
                String resultUrl = "";
                String originalContentUrl = contentUrl;
                System.out.println("===");
                System.out.println("start process combo_id: " + comboId);

                Steps st = cl.PhotolabStepsAdvanced(comboId);
                for (Step v : st.steps) {
                    templateName = Long.toString(v.id);
                    ArrayList<ImageRequest> a = new ArrayList<ImageRequest>();
                    for (String imageUrl : v.image_urls) {
                        if (imageUrl.length() == 0) {
                            imageUrl = originalContentUrl;
                        }
                        ImageRequest im = new ImageRequest();
                        im.url = imageUrl;
                        im.crop = "0,0,1,1";
                        im.flip = 0;
                        im.rotate = 0;
                        a.add(im);
                    }
                    if (a.size() == 0) {
                        ImageRequest im = new ImageRequest();
                        im.url = originalContentUrl;
                        im.crop = "0,0,1,1";
                        im.flip = 0;
                        im.rotate = 0;
                        a.add(im);
                    }
                    ImageRequest[] arr = new ImageRequest[a.size()];
                    arr = a.toArray(arr);
                    resultUrl = cl.TemplateProcess(templateName, arr);
                    originalContentUrl = resultUrl;
                    System.out.println(String.format("--- for template_name: %s, result_url: %s", templateName, resultUrl));
                }
                System.out.println("result_url: " + resultUrl);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
