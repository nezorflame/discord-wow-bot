package main_test

import (
	"log"
	"os"
	"testing"
)

func createMembersList() MembersList {
	var membersList MembersList
	membersList = append(membersList, GuildMember{
		Member: Character{
			Name:  "AAAAA",
			Level: 110,
			Spec: Specialization{
				Name: "Spec1",
			},
			Class: "Class1",
			Items: Items{
				AvgItemLvlEq: 845,
			},
		},
		Rank: 1,
	})
	membersList = append(membersList, GuildMember{
		Member: Character{
			Name:  "BBBB",
			Level: 110,
			Spec: Specialization{
				Name: "Spec2",
			},
			Class: "Class1",
			Items: Items{
				AvgItemLvlEq: 845,
			},
		},
		Rank: 1,
	})
	membersList = append(membersList, GuildMember{
		Member: Character{
			Name:  "CCC",
			Level: 110,
			Spec: Specialization{
				Name: "Spec2",
			},
			Class: "Class1",
			Items: Items{
				AvgItemLvlEq: 847,
			},
		},
		Rank: 1,
	})
	membersList = append(membersList, GuildMember{
		Member: Character{
			Name:  "DD",
			Level: 109,
			Spec: Specialization{
				Name: "Spec1",
			},
			Class: "Class2",
			Items: Items{
				AvgItemLvlEq: 847,
			},
		},
		Rank: 1,
	})
	membersList = append(membersList, GuildMember{
		Member: Character{
			Name:  "E",
			Level: 109,
			Spec: Specialization{
				Name: "Spec1",
			},
			Class: "Class3",
			Items: Items{
				AvgItemLvlEq: 740,
			},
		},
		Rank: 1,
	})
	return membersList
}

func TestStringSorting(t *testing.T) {
	Logger = log.New(os.Stderr, "  ", log.Ldate|log.Ltime)
	members := createMembersList()
	newMembers := members.SortGuildMembers([]string{"name=desc", "level=desc", "ilvl=desc"})
	// for _, m := range newMembers {
	//     log.Println(m.Member.Name, m.Member.Level, m.Member.Class, m.Member.Spec, m.Member.Items.AvgItemLvlEq)
	// }
	if len(newMembers) != len(members) {
		log.Println(newMembers)
		t.Fail()
	}
}
