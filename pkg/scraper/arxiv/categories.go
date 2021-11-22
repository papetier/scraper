package arxiv

import (
	"github.com/gocolly/colly/v2"
	"github.com/papetier/scraper/pkg/database"
	"github.com/papetier/scraper/pkg/scraper/collector"
	log "github.com/sirupsen/logrus"
	"regexp"
)

const (
	arxivCategoryTaxonomyUrl = "https://arxiv.org/category_taxonomy"
	arxivPattern             = `^(.+)\((.+)\)$`
)

var categoriesByCodeMap map[string]*database.ArxivCategory

func UpdateAndLoadCategories(website *database.Website) error {
	log.Info("loading categories")

	// prepare collector
	wc := collector.GetWebsiteCollector(website, colly.AllowURLRevisit())
	wc.Collector.OnHTML("#category_taxonomy_list", categoriesParser)

	// read the categories
	wc.AddUrl(arxivCategoryTaxonomyUrl)
	return nil
}

func categoriesParser(e *colly.HTMLElement) {
	// get the groups
	groupNameList := e.ChildTexts("h2")
	var arxivGroupList []*database.ArxivGroup
	for _, groupName := range groupNameList {
		arxivGroup := &database.ArxivGroup{
			OriginalGroupName: groupName,
		}
		arxivGroupList = append(arxivGroupList, arxivGroup)
	}

	arxivNameRegexp := regexp.MustCompile(arxivPattern)

	// group level
	e.ForEach(".accordion-body", func(groupIndex int, groupBlock *colly.HTMLElement) {
		var arxivArchiveList []*database.ArxivArchive

		// archive level
		groupBlock.ForEach(".accordion-body > .columns", func(archiveIndex int, archiveBlock *colly.HTMLElement) {
			arxivArchive := &database.ArxivArchive{
				ArxivGroup: arxivGroupList[groupIndex],
			}

			archiveFullName := archiveBlock.ChildText("h3")
			if archiveFullName != "" {
				result := arxivNameRegexp.FindStringSubmatch(archiveFullName)
				if len(result) < 3 {
					log.Fatalf("error when extracting the archive's code and name from `%s`", archiveFullName)
				}
				archiveName := result[1]
				archiveCode := result[2]
				arxivArchive.OriginalArxivArchiveName = archiveName
				arxivArchive.OriginalArxivArchiveCode = archiveCode
			}

			var arxivCategoryList []*database.ArxivCategory

			// category level
			archiveBlock.ForEach(".columns", func(categoryIndex int, categoryBlock *colly.HTMLElement) {

				arxivCategory := &database.ArxivCategory{
					ArxivArchive: arxivArchive,
				}

				// category code + name
				categoryFullName := categoryBlock.ChildText("h4")
				if categoryFullName != "" {
					result := arxivNameRegexp.FindStringSubmatch(categoryFullName)
					if len(result) < 3 {
						log.Fatalf("error when extracting the category's code and name from `%s`", categoryFullName)
					}
					categoryName := result[2]
					categoryCode := result[1]
					arxivCategory.OriginalArxivCategoryCode = categoryCode
					arxivCategory.OriginalArxivCategoryName = categoryName
				}

				// category description
				categoryDescription := categoryBlock.ChildText("p")
				arxivCategory.OriginalArxivCategoryDescription = categoryDescription

				// add category to category list
				arxivCategoryList = append(arxivCategoryList, arxivCategory)
			})

			// attach category list to archive
			arxivArchive.ArxivCategories = arxivCategoryList

			// add archive to archive list
			arxivArchiveList = append(arxivArchiveList, arxivArchive)
		})

		// attach archive list to the archive group
		arxivGroupList[groupIndex].ArxivArchives = arxivArchiveList
	})

	// build the category map + list
	categoriesByCodeMap = make(map[string]*database.ArxivCategory)
	for _, group := range arxivGroupList {
		for _, archive := range group.ArxivArchives {
			for _, category := range archive.ArxivCategories {
				categoriesByCodeMap[category.OriginalArxivCategoryCode] = category
			}
		}
	}

	// save the categories in db
	err := database.SaveArxivGroupsArchivesAndCategories(arxivGroupList)
	if err != nil {
		log.Fatalf("saving the arXiv's categories: %s", err)
	}

	log.Info("arXiv categories updated")
}
